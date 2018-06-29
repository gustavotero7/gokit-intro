package main

import (
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {

	logger := log.NewLogfmtLogger(os.Stderr)
	//svc := stringService{}
	svc := loggingMiddleware{
		logger,
		stringService{},
	}

	//uppercaseHandler := httptransport.NewServer(loggingMidleware(log.With(logger, "method", "uppercase"))(makeUppercaseEndpoint(svc)), decodeUppercaseRequest, encodeResponse)
	//countHandler := httptransport.NewServer(loggingMidleware(log.With(logger, "method", "count"))(makeCountEndpoint(svc)), decodeCountRequest, encodeResponse)
	uppercaseHandler := httptransport.NewServer((makeUppercaseEndpoint(svc)), decodeUppercaseRequest, encodeResponse)
	countHandler := httptransport.NewServer((makeCountEndpoint(svc)), decodeCountRequest, encodeResponse)
	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
	http.ListenAndServe(":8080", nil)
}
