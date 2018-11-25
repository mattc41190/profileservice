package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	profileservice "github.com/mattc41190/profileservice"
)

// The entrypoint to our application
func main() {

	// Create a command line flag that will allow users to pass a port that they want the service to run on
	httpAddr := flag.String("http.addr", ":8080", "The port your service should listen on")

	// Parse all flags, this will allow all `flag` var values above to be accessible via their pointer "readthrough" e.g. `*someVar`
	flag.Parse()

	// Create a variable called logger which will hold a value of type `log.Logger`
	var logger log.Logger

	// Create a closure block so that the logger can be safely manipulated. This is unneeded in this file, but is a common pracgtice to
	// avoid confusing var shadowing situations
	{
		// Set the logger value equal to the return of `NewLogfmtLogger`
		// In this case that means all `logger.Log()` calls will end up in Stderr
		logger = log.NewLogfmtLogger(os.Stderr)

		// Prepend all calls to `logger.Log()` with a timestamp prefixed with the value `ts`
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)

		// Prepend all calls to `logger.Log()` with the name of and line of the function which the logger is being called from.
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// Create a variable capable of handling values of type `profileservice.Service`
	var s profileservice.Service

	{
		// Create a new inMemSvc, which is an implementer of the Service interface
		s = profileservice.NewInmemService()

		// Decorate the service with Logging Middleware. Redeclare each function in
		// the Service and do some pre / post(`defer`) processing on the func
		s = profileservice.LoggingMiddleware(logger)(s)
	}

	// Create a value capable of holding values of type `http.Handler`
	var h http.Handler

	{
		// Set the http.Handler `h` to the value returned by `MakeHTTPHandler`
		// Pass the Service implementer `s` and a composoed prefixed logger
		h = profileservice.MakeHTTPHandler(s, log.With(logger, "component", "HTTP"))
	}

	// Make a channel which will operate on error values
	errs := make(chan error)

	// Create and exectue a go routine
	go func() {

		// Make a channel which listens for signals being sent to the program
		c := make(chan os.Signal)

		// When a SIGINT or SIGTERM signal is sent intercept it and add it to the channel
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

		// Push the value from the signal channel to the errs channel
		errs <- fmt.Errorf("%s", <-c)
	}()

	// Create and exectue a go routine
	go func() {

		// Log out the transport protocol being used and the http address the service is using
		logger.Log("transport", "HTTP", "addr", *httpAddr)

		// Send any errors returned from `ListenAndServe` to the `errs` channel
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	// Log `exit` followed by values taken from the error channel.
	// Values may come from the FIRST to come channel. Program will only error and exit once.
	// Note: this line causes the TRUE `main` func to exit.
	logger.Log("exit", <-errs)

}
