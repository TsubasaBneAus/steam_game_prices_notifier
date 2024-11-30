package service

import "net/http"

//go:generate mockgen -source=./httpclient.go -destination=../external/httpclient/mock/httpclient.go -package=mock -typed

// An interface for a HTTP client to send requests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
