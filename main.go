package main

import (
	"context"
	"net/http"
	"os"

	"github.com/go-kit/kit/endpoint"

	"github.com/go-kit/kit/log"

	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {

	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "Listen", ":8080", "called", log.DefaultCaller)

	// Creates an instance of stringService wich implements StringService interface
	svc := stringService{}

	uppercaseHandler := httptransport.NewServer(makeUppercaseEndpoint(svc), decodeUppercaseRequest, encodeResponse)
	countHandler := httptransport.NewServer(makeCountEndpoint(svc), decodeCountRequest, encodeResponse)

	// Here we just define the routes for the cretaed handlers
	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
	http.ListenAndServe(":8080", nil)

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
