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
	// copy headers
	copyheaders(with.header, w.Header())
	// write response
	if err := Write(w, r, with.Code, with.Data); err != nil {
		Err(w, r, &with, err)
	}
}

// SetHeader specifies a response header.
// Headers set this way will overwrite existing headers.
// See http.Header.Set.
func (with With) SetHeader(key, value string) *With {
	with.initheaders()
	with.header.Set(key, value)
	return &with
}

// AddHeader specifies a response header.
// Headers set this way will append to existing headers.
// See http.Header.Add.
func (with With) AddHeader(key, value string) *With {
	with.initheaders()
	with.header.Add(key, value)
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
