package server

import (
	"context"
	"encoding/base64"
	_ "fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	ctx := context.Background()
	RegisterServer(ctx, "urlfunction", NewLambdaURLFunctionServer)
}

// LambdaURLFunctionServer implements the `Server` interface for a use in a AWS LambdaURLFunction + API Gateway context.
type LambdaURLFunctionServer struct {
	Server
	handler http.Handler
}

// NewLambdaURLFunctionServer returns a new `LambdaURLFunctionServer` instance configured by 'uri' which is
// expected to be defined in the form of:
//
//	urlfunction://
func NewLambdaURLFunctionServer(ctx context.Context, uri string) (Server, error) {

	server := LambdaURLFunctionServer{}
	return &server, nil
}

// Address returns the fully-qualified URL used to instantiate 's'.
func (s *LambdaURLFunctionServer) Address() string {
	return ""
}

// ListenAndServe starts the serve and listens for requests using 'mux' for routing.
func (s *LambdaURLFunctionServer) ListenAndServe(ctx context.Context, mux http.Handler) error {
	s.handler = mux
	lambda.Start(s.handleRequest)
	return nil
}

func (s *LambdaURLFunctionServer) handleRequest(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	req, err := newHTTPRequest(ctx, request)

	if err != nil {
		return events.LambdaFunctionURLResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	rec := httptest.NewRecorder()

	s.handler.ServeHTTP(rec, req)

	rsp := rec.Result()

	return events.LambdaFunctionURLResponse{Body: rec.Body.String(), StatusCode: rsp.StatusCode}, nil
}

func newHTTPRequest(ctx context.Context, event events.LambdaFunctionURLRequest) (*http.Request, error) {

	// https://pkg.go.dev/github.com/aws/aws-lambda-go/events#LambdaFunctionURLRequest
	// https://pkg.go.dev/github.com/aws/aws-lambda-go/events#LambdaFunctionURLRequestContextHTTPDescription

	rawQuery := event.RawQueryString

	if len(rawQuery) == 0 {
		params := url.Values{}
		for k, v := range event.QueryStringParameters {
			params.Set(k, v)
		}

		/*
			for k, vals := range event.MultiValueQueryStringParameters {
				params[k] = vals
			}
		*/

		rawQuery = params.Encode()
	}

	// https://docs.aws.amazon.com/apigateway/latest/developerguide/set-up-lambda-proxy-integrations.html
	// If you specify values for both headers and multiValueHeaders, API Gateway V1 merges them into a single list.
	// If the same key-value pair is specified in both, only the values from multiValueHeaders will appear
	// in the merged list.
	headers := make(http.Header)
	for k, v := range event.Headers {
		headers.Set(k, v)
	}

	/*
		for k, vals := range event.MultiValueHeaders {
			headers[http.CanonicalHeaderKey(k)] = vals
		}
	*/

	unescapedPath, err := url.PathUnescape(event.RawPath)
	if err != nil {
		return nil, err
	}
	u := url.URL{
		Host:     headers.Get("Host"),
		Path:     unescapedPath,
		RawQuery: rawQuery,
	}

	// Handle base64 encoded body.
	var body io.Reader = strings.NewReader(event.Body)
	if event.IsBase64Encoded {
		body = base64.NewDecoder(base64.StdEncoding, body)
	}

	req_context := event.RequestContext

	// Create a new request.
	r, err := http.NewRequestWithContext(ctx, req_context.HTTP.Method, u.String(), body)
	if err != nil {
		return nil, err
	}

	// Set remote IP address.
	r.RemoteAddr = req_context.HTTP.SourceIP

	// Set request URI
	r.RequestURI = u.RequestURI()

	r.Header = headers

	return r, nil
}
