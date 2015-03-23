package respond

import (
	"encoding/json"
	"log"
	"net/http"
)

// With describes a response.
type With struct {
	Code   int
	Data   interface{}
	header http.Header
}

// To writes the repsonse.
func (with With) To(w http.ResponseWriter, r *http.Request) {
	// setup headers
	h := w.Header()
	setheaders(Headers, h)
	setheaders(with.header, h)
	// write response
	if err := Write(w, r, with.Code, with.Data); err != nil {
		Err(w, r, &with, err)
	}
}

// Header specifies a response header.
// Headers set this way will overwrite any global headers set
// via respond.Headers.
func (with With) Header(key, value string) *With {
	if with.header == nil {
		with.header = make(http.Header)
	}
	with.header.Set(key, value)
	return &with
}

// Write is the function that actually writes the response.
var Write = func(w http.ResponseWriter, r *http.Request, status int, data interface{}) error {
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// Err is called when an internal error occurs while responding.
var Err = func(w http.ResponseWriter, r *http.Request, with *With, err error) {
	log.Println()
}

// Headers are the http.Header items that will be set on every
// response.
// Use respond.With{}.Header() for response specific headers.
var Headers = make(http.Header)

func setheaders(from, to http.Header) {
	if len(from) == 0 {
		return
	}
	for k, vs := range from {
		for _, v := range vs {
			to.Set(k, v)
		}
	}
}
