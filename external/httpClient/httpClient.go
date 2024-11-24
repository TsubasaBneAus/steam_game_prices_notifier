package httpClient

import "net/http"

//go:generate mockgen -source=httpClient.go -destination=./mock/httpClient_mock.go -package=mock

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type httpClient struct{}

var _ HttpClient = (*httpClient)(nil)

func NewHttpClient() HttpClient {
	return &httpClient{}
}

// Send a HTTP request
func (c *httpClient) Do(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}
