package respond

import "net/http"

func getResponder(r *http.Request) *Responder {
	// get the responder for this request
	mutex.RLock()
	responder, ok := responders[r]
	mutex.RUnlock()
	if !ok {
		panic("respond: must wrap with Handler or HandlerFunc")
	}
	if ManyResponsesPanic && responder == nil {
		// there was a responder there - but it was set to nil - which
		// means we've already responded.
		panic("respond: multiple responses")
	}
	return responder
}

// With writes a response.
func With(w http.ResponseWriter, r *http.Request, status int, data interface{}) {

	responder := getResponder(r)

	// mark the responders[r] as nil - which means we have
	// responded.
	defer func() {
		mutex.Lock()
		responders[r] = nil
		mutex.Unlock()
	}()

	if responder.Before != nil {
		status, data = responder.Before(w, r, status, data)
	}
	if responder.After != nil {
		defer responder.After(w, r, status, data)
	}

	// write the response
	w.WriteHeader(status)
	encoder := responder.Encoder(w, r)
	w.Header().Set("Content-Type", encoder.ContentType(w, r))
	if err := encoder.Encode(w, r, data); err != nil {
		responder.OnErr(w, r, err)
	}

}

// WithStatus responds with the specified status.
// The body will be taken from StatusData.
func WithStatus(w http.ResponseWriter, r *http.Request, status int) {
	responder := getResponder(r)
	With(w, r, status, responder.StatusData(w, r, status))
}

// WithRedirectTemporary sets the Location header and responds with
// the http.StatusTemporaryRedirect status.
func WithRedirectTemporary(w http.ResponseWriter, r *http.Request, location string) {
	responder := getResponder(r)
	w.Header().Set("Location", location)
	With(w, r, http.StatusTemporaryRedirect, responder.RedirectData(w, r, http.StatusTemporaryRedirect, location))
}

// WithRedirectPermanent sets the Location header and responds with
// the http.StatusMovedPermanently status.
func WithRedirectPermanent(w http.ResponseWriter, r *http.Request, location string) {
	responder := getResponder(r)
	w.Header().Set("Location", location)
	With(w, r, http.StatusMovedPermanently, responder.RedirectData(w, r, http.StatusMovedPermanently, location))
}
