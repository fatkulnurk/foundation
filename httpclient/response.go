package httpclient

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Response wraps HTTP response dengan helper methods
type Response struct {
	StatusCode  int
	Body        []byte
	Headers     http.Header
	RawResponse *http.Response
}

// JSON unmarshal response body ke target interface
func (r *Response) JSON(target interface{}) error {
	if len(r.Body) == 0 {
		return fmt.Errorf("empty response body")
	}
	return json.Unmarshal(r.Body, target)
}

// String returns response body as string
func (r *Response) String() string {
	return string(r.Body)
}

// IsSuccess checks if status code is 2xx
func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}
