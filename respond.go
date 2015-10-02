package respond

import (
	"log"
	"net/http"

	"github.com/nbio/httpcontext"
)

var (
	OptionsKey   = "options"
	RespondedKey = "responded"
)

func with(w http.ResponseWriter, r *http.Request, status int, data interface{}, opts *Options, multiple bool) {
	hasOpts := opts != nil

	if hasOpts && multiple && !opts.AllowMultiple {
		panic("respond: multiple responses")
	}

	encoder := JSON // JSON by default
	if hasOpts {
		if opts.Before != nil {
			status, data = opts.Before(w, r, status, data)
		}
		if opts.Encoder != nil {
			encoder = opts.Encoder(w, r)
		}
	}

	// write response
	w.Header().Set("Content-Type", encoder.ContentType(w, r))
	w.WriteHeader(status)
	if err := encoder.Encode(w, r, data); err != nil {
		if hasOpts && opts.OnErr != nil {
			opts.OnErr(err)
		} else {
			panic("respond: " + err.Error())
		}
	}

	if hasOpts {
		if opts.After != nil {
			opts.After(w, r, status, data)
		}
		httpcontext.Set(r, RespondedKey, true)
	}

}

// With responds to the client.
func With(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	var opts *Options
	optsInterface, hasOpts := httpcontext.GetOk(r, OptionsKey)
	if hasOpts {
		opts = optsInterface.(*Options)
	}
	_, multiple := httpcontext.GetOk(r, RespondedKey)
	with(w, r, status, data, opts, multiple)
}

// WithStatus responds to the client with the specified status.
// Options.StatusData will be called to obtain the data payload, or a default
// payload will be returned:
//     {"status":"I'm a teapot","code":418}
func WithStatus(w http.ResponseWriter, r *http.Request, status int) {
	var opts *Options
	optsInterface, hasOpts := httpcontext.GetOk(r, OptionsKey)
	if hasOpts {
		opts = optsInterface.(*Options)
	}
	_, multiple := httpcontext.GetOk(r, RespondedKey)
	var data interface{}
	if hasOpts && opts.StatusData != nil {
		data = opts.StatusData(w, r, status)
	} else {
		const (
			fieldStatus = "status"
			fieldCode   = "code"
		)
		data = map[string]interface{}{fieldStatus: http.StatusText(status), fieldCode: status}
	}
	with(w, r, status, data, opts, multiple)
}

// Options provides additional control over the behaviour of With.
type Options struct {
	// AllowMultiple indicates that multiple responses are allowed. Otherwise,
	// multiple calls to With will panic.
	AllowMultiple bool

	// OnErr is a function field that gets called when an
	// error occurs while responding.
	// By default, the error panic but you may
	// use Options.OnErrLog to just log the error out instead,
	// or provide your own.
	OnErr func(err error)

	// Encoder is a function field that gets the encoder to
	// use to respond to the specified http.Request.
	// If nil, JSON will be used.
	Encoder func(w http.ResponseWriter, r *http.Request) Encoder

	// Before is called for before each response is written
	// and gives user code the chance to mutate the status or data.
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
}

// Handler wraps an HTTP handler becoming the source of options for all
// containing With calls.
func (o *Options) Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpcontext.Set(r, OptionsKey, o)
		handler.ServeHTTP(w, r)
	})
}

func (o *Options) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	httpcontext.Set(r, OptionsKey, o)
}

// OnErrLog prints a log out with the specified error.
// It is an option for Options.OnErr.
func (o *Options) OnErrLog(err error) {
	log.Println("respond: " + err.Error())
}
