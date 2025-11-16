package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// NewRequest creates a new chainable request builder
func (c *client) NewRequest() *Request {
	return &Request{
		client:  c,
		headers: make(map[string]string),
		ctx:     context.Background(),
	}
}

// Request is a chainable request builder
type Request struct {
	client      *client
	ctx         context.Context
	method      string
	url         string
	headers     map[string]string
	body        interface{}
	contentType string
	formData    map[string]string
	formFiles   map[string][]byte
	rawBody     []byte
}

// WithContext sets context for request
func (r *Request) WithContext(ctx context.Context) *Request {
	r.ctx = ctx
	return r
}

// WithMethod sets HTTP method
func (r *Request) WithMethod(method string) *Request {
	r.method = method
	return r
}

// WithURL sets request URL
func (r *Request) WithURL(url string) *Request {
	r.url = url
	return r
}

// WithHeader sets a single header
func (r *Request) WithHeader(key, value string) *Request {
	r.headers[key] = value
	return r
}

// WithHeaders sets multiple headers
func (r *Request) WithHeaders(headers map[string]string) *Request {
	for k, v := range headers {
		r.headers[k] = v
	}
	return r
}

// WithJSON sets body as JSON and content type
func (r *Request) WithJSON(body interface{}) *Request {
	r.body = body
	r.contentType = ContentTypeJSON
	return r
}

// WithFormURLEncoded sets form data as x-www-form-urlencoded
func (r *Request) WithFormURLEncoded(data map[string]string) *Request {
	r.formData = data
	r.contentType = ContentTypeFormURLEncoded
	return r
}

// WithMultipartForm sets form data as multipart/form-data
func (r *Request) WithMultipartForm(data map[string]string, files map[string][]byte) *Request {
	r.formData = data
	r.formFiles = files
	r.contentType = ContentTypeMultipartForm
	return r
}

// WithRaw sets raw body bytes
func (r *Request) WithRaw(body []byte, contentType string) *Request {
	r.rawBody = body
	r.contentType = contentType
	return r
}

// WithText sets plain text body
func (r *Request) WithText(text string) *Request {
	r.rawBody = []byte(text)
	r.contentType = ContentTypeText
	return r
}

// Send executes the request
func (r *Request) Send() (*Response, error) {
	if r.client.config.BaseURL != "" && !strings.HasPrefix(r.url, "http") {
		r.url = r.client.config.BaseURL + r.url
	}

	var lastErr error
	attempts := r.client.config.RetryCount + 1

	for i := 0; i < attempts; i++ {
		if i > 0 {
			select {
			case <-r.ctx.Done():
				return nil, r.ctx.Err()
			case <-time.After(r.client.config.RetryWaitTime):
			}
		}

		resp, err := r.execute()
		if err == nil {
			return resp, nil
		}

		lastErr = err
	}

	return nil, fmt.Errorf("request failed after %d attempts: %w", attempts, lastErr)
}

// execute performs the actual HTTP request
func (r *Request) execute() (*Response, error) {
	var bodyReader io.Reader
	var err error

	// Build body based on content type
	switch r.contentType {
	case ContentTypeJSON:
		if r.body != nil {
			jsonData, err := json.Marshal(r.body)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal JSON: %w", err)
			}
			bodyReader = bytes.NewBuffer(jsonData)
		}

	case ContentTypeFormURLEncoded:
		if r.formData != nil {
			form := url.Values{}
			for k, v := range r.formData {
				form.Set(k, v)
			}
			bodyReader = strings.NewReader(form.Encode())
		}

	case ContentTypeMultipartForm:
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Add form fields
		for k, v := range r.formData {
			if err := writer.WriteField(k, v); err != nil {
				return nil, fmt.Errorf("failed to write form field: %w", err)
			}
		}

		// Add files
		for filename, fileData := range r.formFiles {
			part, err := writer.CreateFormFile("file", filename)
			if err != nil {
				return nil, fmt.Errorf("failed to create form file: %w", err)
			}
			if _, err := part.Write(fileData); err != nil {
				return nil, fmt.Errorf("failed to write file data: %w", err)
			}
		}

		if err := writer.Close(); err != nil {
			return nil, fmt.Errorf("failed to close multipart writer: %w", err)
		}

		bodyReader = body
		r.contentType = writer.FormDataContentType()

	default:
		if len(r.rawBody) > 0 {
			bodyReader = bytes.NewBuffer(r.rawBody)
		}
	}

	req, err := http.NewRequestWithContext(r.ctx, r.method, r.url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set default headers from client config
	for k, v := range r.client.config.DefaultHeaders {
		req.Header.Set(k, v)
	}

	// Set request headers
	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	// Set Content-Type if not already set
	if r.contentType != "" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", r.contentType)
	}

	resp, err := r.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &Response{
		StatusCode:  resp.StatusCode,
		Body:        respBody,
		Headers:     resp.Header,
		RawResponse: resp,
	}, nil
}
