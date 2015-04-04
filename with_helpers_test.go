package respond_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cheekybits/is"
	"github.com/matryer/respond"
)

func TestWithNotFound(t *testing.T) {
	is := is.New(t)
	w := httptest.NewRecorder()
	r := request()
	respond.WithNotFound().To(w, r)
	is.Equal(http.StatusNotFound, w.Code)
	is.Equal(w.Body.String(), `{"error":"not found"}`+"\n")
}

func TestWithRedirectTemporary(t *testing.T) {
	is := is.New(t)
	w := httptest.NewRecorder()
	r := request()
	respond.WithRedirectTemporary("/location").To(w, r)
	is.Equal(http.StatusTemporaryRedirect, w.Code)
	is.Equal(w.Body.String(), `{"redirect-to":"/location"}`+"\n")
	is.Equal(w.Header().Get("Location"), "/location")
}

func TestWithRedirectPermanent(t *testing.T) {
	is := is.New(t)
	w := httptest.NewRecorder()
	r := request()
	respond.WithRedirectPermanent("/location").To(w, r)
	is.Equal(http.StatusMovedPermanently, w.Code)
	is.Equal(w.Body.String(), `{"redirect-to":"/location"}`+"\n")
	is.Equal(w.Header().Get("Location"), "/location")
}

func TestWithErr(t *testing.T) {
	is := is.New(t)
	w := httptest.NewRecorder()
	r := request()
	respond.WithErr(http.StatusUnauthorized).To(w, r)
	is.Equal(http.StatusUnauthorized, w.Code)
	is.Equal(w.Body.String(), `{"error":"Unauthorized"}`+"\n")
}
