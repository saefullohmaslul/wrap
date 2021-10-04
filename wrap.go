package wrap

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	wrap "github.com/vdparikh/wrap/v4"

	"github.com/aws/aws-lambda-go/events"
	"github.com/labstack/echo"
)

// Route wraps echo server into Lambda Handler
func Route(e *echo.Echo) func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		body := strings.NewReader(request.Body)
		req := httptest.NewRequest(request.HTTPMethod, request.Path, body)
		for k, v := range request.Headers {
			req.Header.Add(k, v)
		}

		q := req.URL.Query()
		for k, v := range request.QueryStringParameters {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		res := rec.Result()
		responseBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return wrap.FormatAPIResponse(http.StatusInternalServerError, res.Header, err.Error())
		}

		return wrap.FormatAPIResponse(res.StatusCode, res.Header, string(responseBody))
	}
}
