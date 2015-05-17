package respond_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cheekybits/is"
	"github.com/matryer/respond"
)

func TestWithStatus(t *testing.T) {
	is := is.New(t)

	responder := respond.New()
	w := httptest.NewRecorder()
	r := newTestRequest()
	responder.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respond.WithStatus(w, r, http.StatusTeapot)
	})(w, r)

	// assert it was written
	is.Equal(http.StatusTeapot, w.Code)
	var data map[string]interface{}
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &data))
	is.Equal(data["status"], http.StatusText(http.StatusTeapot))
	is.Equal(data["code"], http.StatusTeapot)

}

func TestWithRedirectTemporary(t *testing.T) {
	is := is.New(t)

	responder := respond.New()
	w := httptest.NewRecorder()
	r := newTestRequest()
	responder.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respond.WithRedirectTemporary(w, r, "/new/path")
	})(w, r)

	// assert it was written
	is.Equal(http.StatusTemporaryRedirect, w.Code)
	var data map[string]interface{}
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &data))
	is.Equal(data["status"], http.StatusText(http.StatusTemporaryRedirect))
	is.Equal(data["code"], http.StatusTemporaryRedirect)
	is.Equal(w.HeaderMap.Get("Location"), "/new/path")

}
