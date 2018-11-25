package client

import (
	"io"
	"log"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"

	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	consulapi "github.com/hashicorp/consul/api"

	profileservice "github.com/mattc41190/go-learn-go/gokit-examples/02-profileservice"
)

// TODO:
// A. Understand how consul works.
// B. Understand how the consul API works

// New is a function which accepts a string and a Logger and returns an implementer of `profileservice.Service` and `error`
func New(consulAddr string, logger log.Logger) (profileservice.Service, error) {

	// Create a variable called `endpoints` capable of holding `profileservice.Endpoints`
	var endpoints profileservice.Endpoints

	// consulService is the name of the service in Consul
	consulService := "profileservice"

	// consulTags is a list of Consul filter tags
	consulTags := []string{"prod"}

	// passingOnly is a shorthand for a filtering boolean
	// If passingOnly is true, only instances where both the service and any proxy are healthy will be returned
	passingOnly := true

	// retryMax is the max number of times to attempt to get a consul service showing "healthy"
	retryMax := 4

	// retryTimeout is the amount of time before a health check times out
	retryTimeout := time.Millisecond * 500

	// Create a consulapi.Config struct and set it to a variable named `apiConfig`
	apiConfig := consulapi.Config{Address: consulAddr}

	// Create variable named `apiClient` and `err` and set their values to the return values from `consulapi.NewClient`
	// A pointer to a value of type `consulapi.Config` (`apiConfig` in this case) is passed to the NewClient code.
	apiClient, err := consulapi.NewClient(&apiConfig)

	// Handle any errors
	if err != nil {
		logger.Log("error", err)
		return nil, err
	}

	sdclient := consul.NewClient(apiClient)
	instancer := consul.NewInstancer(sdclient, logger, consulService, consulTags, passingOnly)

	// The following section is incredibly complicated... Get Ready
	// Create a new closure to avoid variable shadowing
	{
		// Create a variable called factory which accepts a function `profileservice.MakePostProfileEndpoint`
		// and returns a function (sd.Factory) which accepts a string an returns an `Endpoint` function (among other things)
		// One SHOULD see `factory` as a function.
		// REWORDING FOR MY HEAD
		// Create a variable called `factory` whose value is an `sd.Factory` created by calling `factoryFor` on an function of type Endpoint.
		// Note: `sd.Factory` is a type function and thusly `factory` is a function. This continually confuses me
		factory := factoryFor(profileservice.MakePostProfileEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.PostProfile = retry
	}
	{
		factory := factoryFor(profileservice.GetProfileEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.PostProfile = retry
	}
	{
		factory := factoryFor(profileservice.PutProfileEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.PostProfile = retry
	}
	{
		factory := factoryFor(profileservice.PatchProfileEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.PostProfile = retry
	}
	{
		factory := factoryFor(profileservice.DeleteProfileEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.PostProfile = retry
	}
	{
		factory := factoryFor(profileservice.GetAddressesEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.PostProfile = retry
	}
	{
		factory := factoryFor(profileservice.GetAddressEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.PostProfile = retry
	}
	{
		factory := factoryFor(profileservice.PostAddressEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.PostProfile = retry
	}
	{
		factory := factoryFor(profileservice.DeleteAddressEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.PostProfile = retry
	}

	return endpoints

}

// factoryFor will accept a function in this case a named function called `makeEndpoint`
// It will return an sd.Factory (`sd` == "Service Discovery")
// makeEndpoint is a function which accepts an implementer of profileservice.Service
// and returns an Endpoint. In essence factoryFor is a function which accepts a function
// which returns a function which returns a function <- No typo...
// The final function of type `sd.Factory` is a function which accepts a string (service discovery address)
// and returns an Endpoint function (already mentioned), an io.Closer implemeter, and an error
// Quote: The function accepts a function which returns a function, and returns a function.
func factoryFor(makeEndpoint func(profileservice.Service) endpoint.Endpoint) sd.Factory {

	// Create a function to return which accepts a string called `instance`
	// and returns: an Endpoint (function), an implementer of io.Closer, and an error (implementer lol)
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {

		// Create variables (`service` and `err`) for the returns values of `MakeClientEndpoints`.
		// Pass the string `instance` to `MakeClientEndpoints`. Recall that MakeClientEndpoints returns an implementer of profileservice.Service
		// To do this it declares functions on the profileservice.Endpoints struct. These methods will in turn create a request and send it to that same
		// struct's corresponding XEndpoint field whose value is a function of type Endpoint. So the Service in this case
		// has a collection of client type endpoints whose structure is determined by GoKit
		// 1 -- Serialize Request 2 -- Make Request 3 -- Deserialize Response 4 -- Return Response.
		// So using the service AFTER it's declaration would be possible in the following fashion:
		// ctx =  Context.ToDo()
		// p, err := service.GetProfile(ctx, "<uuid>")
		service, err := profileservice.MakeClientEndpoints(instance)

		// Check to see if there was an error creating the Client Service
		if err != nil {

			// Return nil for Endpoint, the closer, and return the error value returned from MakeClientEndpoints
			return nil, nil, err
		}

		// This is a what the fuck moment for me...
		// Recall that to the factory function we pass in a named function declaration.
		// Within the this function we now CALL that function and what we expect to get back
		// is the version of the Endpoint which we would have expected to be on the server I thought.
		// Since `sd.Factory` will contain this bad boy we should see how that function intends to call th service
		return makeEndpoint(service), nil, nil
	}
}

// A. What is the consul API client?

// WIP

// So the sd.Factory / `factory` function (which returns an Endpoint function!) is passed to `sd.NewEndpointer`
// NewEndpointer accepts an instancer Instancer, a factory function, and a logger
