package httprouter

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

type Response struct {
	w          http.ResponseWriter
	statusCode int
	headers    http.Header
	wrote      bool
}

// response(w) → *Response
func ResponseOf(w http.ResponseWriter) *Response {
	return &Response{
		w:          w,
		statusCode: http.StatusOK,
		headers:    make(http.Header),
	}
}

// Status() → set HTTP status (chainable)
func (r *Response) Status(code int) *Response {
	r.statusCode = code
	return r
}

// Header() → set header custom (chainable)
func (r *Response) Header(key, value string) *Response {
	r.headers.Set(key, value)
	return r
}

func (r *Response) writeHeaders(contentType string) {
	if r.wrote {
		return
	}
	r.wrote = true

	// set content-type
	if contentType != "" {
		r.headers.Set("Content-Type", contentType)
	}

	// assign semua header ke ResponseWriter
	for k, vals := range r.headers {
		for _, v := range vals {
			r.w.Header().Add(k, v)
		}
	}

	r.w.WriteHeader(r.statusCode)
}

// JSON() → encode nilai jadi JSON
func (r *Response) JSON(v interface{}) {
	r.writeHeaders("application/json; charset=utf-8")
	if v == nil {
		return
	}
	enc := json.NewEncoder(r.w)
	enc.SetEscapeHTML(true)
	_ = enc.Encode(v)
}

// HTML() → kirim konten HTML
func (r *Response) HTML(html string) {
	r.writeHeaders("text/html; charset=utf-8")
	if html == "" {
		return
	}
	_, _ = r.w.Write([]byte(html))
}

// XML() → encode nilai jadi XML
func (r *Response) XML(v interface{}) {
	r.writeHeaders("application/xml; charset=utf-8")
	if v == nil {
		return
	}
	enc := xml.NewEncoder(r.w)
	_ = enc.Encode(v)
}

// Text() → plain text
func (r *Response) Text(s string) {
	r.writeHeaders("text/plain; charset=utf-8")
	if s == "" {
		return
	}
	_, _ = r.w.Write([]byte(s))
}

// WriteJSON mengirim response JSON dengan status code
func WriteJSON(w http.ResponseWriter, statusCode int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if v == nil {
		return
	}

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)

	if err := enc.Encode(v); err != nil {
		// fallback kalau encoding gagal
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// WriteHTML mengirim response HTML (string mentah) dengan status code
func WriteHTML(w http.ResponseWriter, statusCode int, html string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)

	if html == "" {
		return
	}

	_, _ = w.Write([]byte(html))
}

// WriteXML mengirim response XML dengan status code
func WriteXML(w http.ResponseWriter, statusCode int, v interface{}) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(statusCode)

	if v == nil {
		return
	}

	enc := xml.NewEncoder(w)
	if err := enc.Encode(v); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
