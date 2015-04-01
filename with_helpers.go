package respond

import (
	"errors"
	"net/http"
)

const (
	headerLocation = "Location"
)

// ErrNotFound is the error that is responded when WithNotFound
// is called.
var ErrNotFound = errors.New("not found")

// WithNotFound indicates a NotFound response.
func WithNotFound() *W {
	return With(http.StatusNotFound, ErrNotFound)
}

// RedirectResponse is a func that gets the object to respond with during a
// redirection via WithRedirectTemporary or WithRedirectPermanent.
var RedirectResponse = func(location string) interface{} {
	return map[string]interface{}{"redirect-to": location}
}

// WithRedirectTemporary indicates a temporary redirect response.
func WithRedirectTemporary(location string) *W {
	return With(http.StatusTemporaryRedirect, RedirectResponse(location)).
		SetHeader(headerLocation, location)
}

// WithRedirectPermanent indicates a permanent redirect response.
func WithRedirectPermanent(location string) *W {
	return With(http.StatusMovedPermanently, RedirectResponse(location)).
		SetHeader(headerLocation, location)
}
