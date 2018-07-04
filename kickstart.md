## Go-Kit Kickstart guide
Let’s create a minimal Go kit service. Based in the offical go-kit guide/example stringsvc3, this "rewrite" and attempt to explain better some concepts due the original one have "lost info" (missing or not very clear information in some cases)

### Business logic
In go-kit we model a service as an interface. for now all our code will be inside main.go file (yes, I know it's a mess, but later we will refactor our code)
```
//We will be using Context package, let's import it
import "context.Context"

type StringService interface {
	Uppercase(context.Context, string) (string, error)
	Count(context.Context, string) int
}
```

The interface will have an implementation
```
// At this point we need to import some extra packages, our import stament must look like
import (
	"context"
	"errors"
	"strings"
)

type stringService struct{}

// These methods will process the data parsed by the request decoder (keep reading)
func (stringService) Uppercase(ctx context.Context, str string) (string, error) {
	if str == "" {
		return "", errors.New("Empty string")
	}
	return strings.ToUpper(str), nil
}

func (stringService) Count(ctx context.Context, str string) int {
	return len(str)
}
```

### Requests and responses
In Go kit, the primary messaging pattern is RPC. So, every method in our interface will be modeled as a remote procedure call. For each method, we define request and response structs, capturing all of the input and output parameters respectively.

```
// Model expected request body for uppercase procedure
type uppercaseRequest struct {
	Str string `json:"s"`
}

// Model server response for uppercase procedure
type uppercaseResponse struct {
	Result string `json:"result"`
	Error  string `json:"error"` // Errors don't JSON-Marshal, so we use a string
}

// Model expected request body for count procedure
type countRequest struct {
	Str string `json:"str"`
}

// Model server response for count procedure
type countReponse struct {
	Count int `json:"count"`
}
```

### Endpoints
Go kit provides much of its functionality through an abstraction called an endpoint.
`type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)`

An endpoint represents a single RPC. That is, a single method in our service interface. We’ll write simple adapters to convert each of our service’s methods into an endpoint. Each adapter takes a StringService, and returns an endpoint that corresponds to one of the methods.

```
func makeUppercaseEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// Assert request type, here we parse/marshal the request into a uppercaseRequest object
		req := request.(uppercaseRequest)
		// Call uppercase method passing Str parameter from our decoded request object
		str, err := svc.Uppercase(ctx, req.Str)

		if err != nil {
			// Remember we have defided the erro as an string in the  uppercaseResponse struct?
			// To fill the error field we convert it to a string calling Error() method from the err object
			return uppercaseResponse{"", err.Error()}, nil
		}

		// Return a new uppercaseResponse object
		return uppercaseResponse{str, ""}, nil
	}
}

func makeCountEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// Assert request type, here we parse/marshal the request into a uppercaseRequest object
		req := request.(countRequest)
		// Return countReponse object
		return countReponse{svc.Count(ctx, req.Str)}, nil
	}
}
```

### Transports
Now we need to expose your service to the outside world, so it can be called. Your organization probably already has opinions about how services should talk to each other. Maybe you use Thrift, or custom JSON over HTTP. Go kit supports many transports out of the box.

For this minimal example service, let’s use JSON over HTTP. Go kit provides a helper struct, in package transport/http.

```
func main() {
	// Creates an instance of stringService wich implements StringService interface
	svc := stringService{}

	// Each handler should have an endpoint, a decoder, and a encoder
	uppercaseHandler := httptransport.NewServer(makeUppercaseEndpoint(svc), decodeUppercaseRequest, encodeResponse)
	countHandler := httptransport.NewServer(makeCountEndpoint(svc), decodeCountRequest, encodeResponse)

	// Here we just define the routes for the cretaed handlers
	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
	http.ListenAndServe(":8080", nil)

}

// Parse request body into a uppercaseRequest object
func decodeUppercaseRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req uppercaseRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// Parse request body into a countRequest object
func decodeCountRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req countRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// We can use one encode procedure for both calls, since we only need to encode the provided response into a json object
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

```

Now you can try the servie, lets call uppercase and count endpoints
```
# Call
$ curl -XPOST -d'{"s":"hello, world"}' localhost:8080/uppercase
# Response
{"v":"HELLO, WORLD","err":null}

# Call
$ curl -XPOST -d'{"s":"hello, world"}' localhost:8080/count
# Reponse
{"v":12}
```
## The complete code should looks like this

<details><summary>See code</summary>
<p>
<b>main.go</b>

```golang
package main

import (
	"context"
	"net/http"
    "os"
	"encoding/json"
	"errors"
	"strings"

	"github.com/go-kit/kit/endpoint"

	"github.com/go-kit/kit/log"

	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	// Creates an instance of stringService wich implements StringService interface
	svc := stringService{}

	// Each handler should have an endpoint, a decoder, and a encoder
	uppercaseHandler := httptransport.NewServer(makeUppercaseEndpoint(svc), decodeUppercaseRequest, encodeResponse)
	countHandler := httptransport.NewServer(makeCountEndpoint(svc), decodeCountRequest, encodeResponse)

	// Here we just define the routes for the cretaed handlers
	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
	http.ListenAndServe(":8080", nil)

}


// StringService _
type StringService interface {
	Uppercase(context.Context, string) (string, error)
	Count(context.Context, string) int
}

type stringService struct{}

func (stringService) Uppercase(ctx context.Context, str string) (string, error) {
	if str == "" {
		return "", errors.New("Empty string")
	}
	return strings.ToUpper(str), nil
}

func (stringService) Count(ctx context.Context, str string) int {
	return len(str)
}

// Model expected client request for uppercase resource
type uppercaseRequest struct {
	Str string `json:"str"`
}

// Model server response for uppercase resource
type uppercaseResponse struct {
	Result string `json:"result"`
	Error  string `json:"error"` // Errors don't JSON-Marshal, so we use a string
}

// Model expected client request for count resource
type countRequest struct {
	Str string `json:"str"`
}

// Model server response for count resource
type countReponse struct {
	Count int `json:"count"`
}

// Parse request body into a uppercaseRequest object
func decodeUppercaseRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req uppercaseRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// Parse request body into a countRequest object
func decodeCountRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req countRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// We can use one encode procedure for both calls, since we only need to encode the provided response into a json object
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func makeUppercaseEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// Assert request type, here we parse/marshal the request into a uppercaseRequest object
		// For this the decodeUppercaseRequest function is called, we especify the decoder for this endpoint when the handler is created
		req := request.(uppercaseRequest)
		// Call uppercase method passing Str parameter from our decoded request object
		str, err := svc.Uppercase(ctx, req.Str)

		if err != nil {
			// Remember we have defided the erro as an string in the  uppercaseResponse struct?
			// To fill the error field we convert it to a string calling Error() method from the err object
			return uppercaseResponse{"", err.Error()}, nil
		}

		// Return a new uppercaseResponse object
		return uppercaseResponse{str, ""}, nil
	}
}

func makeCountEndpoint(svc StringService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// Assert request type, here we parse/marshal the request into a uppercaseRequest object
		// For this the decodeUppercaseRequest function is called, we especify the decoder for this endpoint when the handler is created
		req := request.(countRequest)
		// Return countReponse object
		return countReponse{svc.Count(ctx, req.Str)}, nil
	}
}

```

</p>
</details>