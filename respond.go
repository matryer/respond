package respond

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

var (
	mutex      sync.RWMutex
	responders = make(map[*http.Request]*Responder)
	responded  = make(map[*http.Request]bool)
)

// ManyResponsesPanic indicates whether repsond will panic if
// user code attempts to make multiple responses.
var ManyResponsesPanic = true

// Responder provides the ability to respond using repsond.With calls.
type Responder struct {
	// Encoder is a function field that gets the encoder to
	// use to respond to the specified http.Request.
	// By default, just returns the JSON Encoder.
	Encoder func(w http.ResponseWriter, r *http.Request) Encoder

	// OnErr is a function field that is called when an error occurs
	// during the responding process.
	// OnErrLog and OnErrPanic are two provided options with
	// OnErrLog being the default.
	OnErr func(w http.ResponseWriter, r *http.Request, err error)

	// Before is called for before each response is written and gives user code the chance
	// to mutate the status or data.
	// Useful for handling different types of data differently (like errors),
	// enveloping the response, setting common HTTP headers etc.
	Before func(w http.ResponseWriter, r *http.Request, status int, data interface{}) (int, interface{})

	// After is called after each response.
	// Useful for logging activity after a response has been written.
	After func(w http.ResponseWriter, r *http.Request, status int, data interface{})

	// StatusData is a function field that gets the data to respond with when
	// WithStatus is called.
	// By default, the function will return an object that looks like this:
	//     {"status":"Not Found","code":404}
	StatusData func(w http.ResponseWriter, r *http.Request, status int) interface{}

	// RedirectData is a function field that gets the data to respond with when
	// WithRedirect* is called.
	// By default, the function will return an object that looks like this:
	//     {"status":"Temporary Redirect","code":301,"location":"/path"}
	RedirectData func(w http.ResponseWriter, r *http.Request, status int, location string) interface{}
}

// DefaultStatusData is the default StatusData function for a Responder.
func DefaultStatusData(w http.ResponseWriter, r *http.Request, status int) interface{} {
	return map[string]interface{}{
		"status": http.StatusText(status),
		"code":   status,
	}
}

// DefaultRedirectData is the default RedirectData function for a Responder.
func DefaultRedirectData(w http.ResponseWriter, r *http.Request, status int, location string) interface{} {
	return map[string]interface{}{
		"status":   http.StatusText(status),
		"code":     status,
		"location": location,
	}
}

// New makes a new Responder.
func New() *Responder {
	return &Responder{
		Encoder:      func(http.ResponseWriter, *http.Request) Encoder { return JSON },
		OnErr:        OnErrLog,
		StatusData:   DefaultStatusData,
		RedirectData: DefaultRedirectData,
	}
}

// OnErrLog logs the error using log.Println.
func OnErrLog(_ http.ResponseWriter, _ *http.Request, err error) {
	log.Println("respond:", err)
}

// OnErrPanic panics with the specified error.
func OnErrPanic(_ http.ResponseWriter, _ *http.Request, err error) {
	panic(fmt.Sprint("respond:", err))
}

// Handler wraps an http.Handler and enables its handlers to use
// respond.With.
func (d *Responder) Handler(handler http.Handler) http.Handler {
	return d.HandlerFunc(handler.ServeHTTP)
}

// HandlerFunc wraps an http.HandlerFunc and enables the handler to
// use repsond.With.
func (d *Responder) HandlerFunc(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		responders[r] = d
		mutex.Unlock()
		defer func() {
			mutex.Lock()
			delete(responders, r)
			delete(responded, r)
			mutex.Unlock()
		}()
		fn(w, r)
	}
}

// Encoder descirbes an object capable of encoding
// a response.
type Encoder interface {
	// Encode writes a serialization of v to w, optionally using additional
	// information from the http.Request to do so.
	Encode(w http.ResponseWriter, r *http.Request, v interface{}) error
	// ContentType gets a string that will become the Content-Type header
	// when responding through w to the specified http.Request.
	// Most of the time the argument will be ignored, but occasionally
	// details in the request, or even in the headers in the ResponseWriter may
	// change the content type.
	ContentType(w http.ResponseWriter, r *http.Request) string
}

type jsonEncoder struct{}

var _ Encoder = (*jsonEncoder)(nil)

// JSON is an Encoder for JSON.
var JSON *jsonEncoder

func (*jsonEncoder) Encode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

func (*jsonEncoder) ContentType(w http.ResponseWriter, r *http.Request) string {
	return "application/json; charset=utf-8"
}
