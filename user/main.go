package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	action := flag.String("do", "get", "Action to perform")
	id := flag.String("id", "001", "Profile ID to retrieve -- Works with `get`")
	filePath := flag.String("filepath", "", "Path to user profile JSON -- Works with `create")

	flag.Parse()

	switch *action {
	case "get":
		get(*id)
	case "create":
		create(*filePath)
	default:
		panic("User must specify valid 'do' command")
	}
}

func create(path string) {
	file, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	contents, err := ioutil.ReadAll(file)

	if err != nil {
		panic(err)
	}

	resp, err := http.Post("http://localhost:8080/profiles/", "application/json", strings.NewReader(string(contents)))

	if err != nil {
		panic(err)
	}

	fmt.Println(resp.StatusCode)
}

func get(id string) {

	resp, err := http.Get("http://localhost:8080/profiles/" + id)

	if err != nil {
		panic(err)
	}

	fmt.Println(resp.StatusCode)

	contents, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(contents))
}
