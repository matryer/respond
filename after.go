package respond

import (
	"bytes"
	"log"
	"net/http"
	"sync"
)

// AfterFunc is a function that can be called after each response.
type AfterFunc func(w *Response, r *http.Request, status int, data interface{})

// LogAfter is an AfterFunc that logs the repsonse detials via log.Println.
// respond.After(respond.LogAfter) provides a debug level view into
// what responses are being made.
func LogAfter(w *Response, r *http.Request, status int, data interface{}) {
	if w.HasBody() {
		log.Println("respond:", status, w.Body().String())
	} else {
		log.Println("respond:", status, data)
	}
}

// After sets the AfterFunc to call after each response.
func After(fn AfterFunc) {
	afterLock.Lock()
	after = fn
	afterLock.Unlock()
}

// KeepBody indicates whether the Response in After will
// make the Body available or not.
// By default, the Body() method will panic, but setting KeepBody(true)
// will cause the response body to be written to an internal buffer
// of the Response, as well as to the client.
func KeepBody(keep bool) {
	afterLock.Lock()
	keepbody = keep
	afterLock.Unlock()
}

var after AfterFunc
var keepbody bool
var afterLock sync.RWMutex

// Response represents the response given the client.
type Response struct {
	w        http.ResponseWriter
	keepbody bool
	body     *bytes.Buffer
	status   int
}

// Header gets the response headers for the underlying
// http.ResponseWriter.
func (r *Response) Header() http.Header {
	return r.w.Header()
}

// WriteHeader writes the headers with specified status code
// to the underlying http.ResponseWriter.
func (r *Response) WriteHeader(status int) {
	r.w.WriteHeader(status)
}

// Write writes the bytes to the underlying http.ResponseWriter
// and if KeepBody(true), to an internal buffer.
func (r *Response) Write(b []byte) (int, error) {
	if r.keepbody {
		r.body.Write(b)
	}
	return r.w.Write(b)
}

// Body gets the bytes that were written to the response.
// Will panic if HasBody() returns false.
// Use respond.KeepBody(true) to keep copies of bodies.
func (r *Response) Body() *bytes.Buffer {
	if !r.keepbody {
		panic("respond: cannot call Body() when KeepBody(false)")
	}
	return r.body
}

// Status gets the HTTP Status Code that repsond replied
// with.
func (r *Response) Status() int {
	return r.status
}

// HasBody gets whether the Response has a copy of the
// Body or not.
func (r *Response) HasBody() bool {
	return r.keepbody
}
