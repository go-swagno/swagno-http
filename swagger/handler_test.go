package swagger_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-swagno/swagno"
	"github.com/go-swagno/swagno-http/swagger"
	"github.com/go-swagno/swagno/components/endpoint"
	"github.com/go-swagno/swagno/components/http/response"
	"github.com/go-swagno/swagno/components/mime"
	"github.com/go-swagno/swagno/example/models"
	"github.com/stretchr/testify/assert"
)

func TestSwaggerHandler(t *testing.T) {
	sw := swagno.New(swagno.Config{Title: "Testing API", Version: "v1.0.0"})
	desc := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed id malesuada lorem, et fermentum sapien. Vivamus non pharetra risus, in efficitur leo. Suspendisse sed metus sit amet mi laoreet imperdiet. Donec aliquam eros eu blandit feugiat. Quisque scelerisque justo ac vehicula bibendum. Fusce suscipit arcu nisl, eu maximus odio consequat quis. Curabitur fermentum eleifend tellus, lobortis hendrerit velit varius vitae."

	endpoints := []*endpoint.EndPoint{
		endpoint.New(
			endpoint.GET,
			"/product",
			endpoint.WithTags("product"),
			endpoint.WithSuccessfulReturns([]response.Response{response.New(models.Product{}, "200", "OK")}),
			endpoint.WithErrors([]response.Response{response.New(models.UnsuccessfulResponse{}, "400", "Bad Request")}),
			endpoint.WithDescription(desc),
			endpoint.WithProduce([]mime.MIME{mime.JSON, mime.XML}),
			endpoint.WithConsume([]mime.MIME{mime.JSON}),
			endpoint.WithSummary("this is a test summary"),
		),
	}

	sw.AddEndpoints(endpoints)
	handler := swagger.SwaggerHandler(sw.MustToJson())

	server := httptest.NewServer(handler)

	status, err := checkUrl(server.URL + "/swagger/index.html")
	assert.NoError(t, err)
	assert.Equal(t, 200, status)

	status, err = checkUrl(server.URL + "/swagger/doc.json")
	assert.NoError(t, err)
	assert.Equal(t, 200, status)
}

func checkUrl(url string) (statusCode int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return resp.StatusCode, err
	}

	return resp.StatusCode, nil
}

func TestSwaggerHandler_MultipleListeners(t *testing.T) {
	sw1 := swagno.New(swagno.Config{Title: "Testing API 1", Version: "v1.0.0"})
	sw1.AddEndpoints([]*endpoint.EndPoint{endpoint.New(endpoint.GET, "/product1")})
	handler1 := swagger.SwaggerHandler(sw1.MustToJson())
	server1 := httptest.NewServer(handler1)

	checkSwaggerTitle(t, server1.URL, sw1.Info.Title)

	sw2 := swagno.New(swagno.Config{Title: "Testing API 2", Version: "v1.0.0"})
	sw2.AddEndpoints([]*endpoint.EndPoint{endpoint.New(endpoint.GET, "/product2")})
	handler2 := swagger.SwaggerHandler(sw2.MustToJson())
	server2 := httptest.NewServer(handler2)

	checkSwaggerTitle(t, server2.URL, sw2.Info.Title)

	assert.NotEqual(t, server1.URL, server2.URL)
}

func checkSwaggerTitle(t *testing.T, serverUrl string, expectedTitle string) {
	resp, err := http.Get(serverUrl + "/swagger/doc.json")
	assert.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.True(t, strings.Contains(string(body), expectedTitle))
}
