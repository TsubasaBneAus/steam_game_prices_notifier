package httpclient

import (
	"net/http"

	"github.com/TsubasaBneAus/steam_game_price_notifier/app/service"
)

type httpClient struct{}

var _ service.HTTPClient = (*httpClient)(nil)

// Generate a new HTTP client
func NewHTTPClient() *httpClient {
	return &httpClient{}
}

// Send a HTTP request
func (c *httpClient) Do(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}
