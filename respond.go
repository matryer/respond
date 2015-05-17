package respond

import (
	"net/http"
	"sync"
)

// Responder handles responses.
type Responder struct {
	// m is the map of Response objects to http.Request objects.
	m map[*http.Request]*Response
	// wrapped keeps track of what has been wrapped or not.
	wrapped map[*http.Request]bool
	// encoders holds the Encoder objects that this Responder
	// will support.
	encoders *encodersMap
	// DefaultEncoder is the Encoder to use when a better match
	// it not found.
	DefaultEncoder Encoder
	// lock is the sync.RWMutex used to ensure safety.
	lock sync.RWMutex
}

// New makes a new Responder.
// A Responder can safely respond to many requests.
// Typically you would have one Responder per app.
func New() *Responder {
	return &Responder{
		m:       make(map[*http.Request]*Response),
		wrapped: make(map[*http.Request]bool),
		encoders: &encodersMap{
			encoders: map[string]Encoder{"json": JSON},
		},
		DefaultEncoder: JSON,
	}
}

// To specifies the http.Request to respond to via the http.ResponseWriter.
func (d *Responder) To(w http.ResponseWriter, r *http.Request) *Response {
	response := &Response{d: d, w: w, r: r}
	d.lock.Lock()
	if !d.wrapped[r] {
		panic("respond: must wrap with Handler or HandlerFunc")
	}
	d.m[r] = response
	d.lock.Unlock()
	return response
}

func (d *Responder) setup(r *http.Request) {
	d.lock.Lock()
	d.wrapped[r] = true
	d.lock.Unlock()
}

func (d *Responder) finish(w http.ResponseWriter, r *http.Request) {
	// get the response if it's there
	d.lock.RLock()
	response, ok := d.m[r]
	d.lock.RUnlock()
	if !ok {
		// no response - skip, nothing to do
		return
	}

	// write the response
	d.write(w, r, response)

	// clean up
	d.lock.Lock()
	delete(d.m, r)
	d.lock.Unlock()
}

func (d *Responder) write(w http.ResponseWriter, r *http.Request, response *Response) {

	// find the encoder
	var encoder Encoder
	var ok bool
	if encoder, ok = d.encoders.Match(r.Header.Get("Accept")); !ok {
		encoder = d.DefaultEncoder
	}

	// write response
	w.Header().Set("Content-Type", encoder.ContentType(w, r))
	w.WriteHeader(response.status)
	encoder.Encode(w, r, response.data)

}

// Encoders allows the adding, removing and matching of encoders
// by fuzzy strings, usually the Accept request header.
// By default, JSON is added.
func (d *Responder) Encoders() interface {
	Add(match string, encoder Encoder)
	Match(s string) (Encoder, bool)
	Del(e Encoder)
} {
	return d.encoders
}

// Handler wraps a http.Handler that makes use of respond.
// Using respond without wrapping the handlers will panic.
func (d *Responder) Handler(handler http.Handler) http.Handler {
	return d.HandlerFunc(handler.ServeHTTP)
}

// HandlerFunc wraps a http.HandlerFunc that makes use of respond.
// Using respond without wrapping the handlers will panic.
func (d *Responder) HandlerFunc(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d.setup(r)
		fn(w, r)
		d.finish(w, r)
	}
}

// Response represents a single response.
type Response struct {
	d *Responder
	w http.ResponseWriter
	r *http.Request

	// response details
	status int
	data   interface{}
}

// With specifies the HTTP status code and any data to respond with.
func (r *Response) With(status int, data interface{}) *Response {
	r.status = status
	r.data = data
	return r
}

// Now writes the response and clears things up immediately rather
// than waiting until the handler is finished.
// Should be the last call on a Response.
func (r *Response) Now() {
	r.d.finish(r.w, r.r)
}
