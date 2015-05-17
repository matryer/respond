package respond_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cheekybits/is"
	"github.com/matryer/respond"
)

var testdata = map[string]interface{}{"test": true}

func request() *http.Request {
	r, _ := http.NewRequest("GET", "Something", nil)
	return r
}

type testHandler struct {
	respond *respond.Responder
}

func (t *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.respond.To(w, r).With(http.StatusOK, testdata)
}

func TestToWithHandler(t *testing.T) {
	is := is.New(t)

	respond := respond.New()
	w := httptest.NewRecorder()
	r := request()

	handler := &testHandler{respond}
	respond.Handler(handler).ServeHTTP(w, r)

	is.Equal(http.StatusOK, w.Code)
	var data map[string]interface{}
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &data))
	is.Equal(data, testdata)
}

// TestTo tests the simple respond.To(w, r).With(status, data) case.
func TestToWithHandlerFunc(t *testing.T) {
	is := is.New(t)

	respond := respond.New()
	w := httptest.NewRecorder()
	r := request()

	respond.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respond.To(w, r).With(http.StatusOK, testdata)
	})(w, r)

	is.Equal(http.StatusOK, w.Code)
	var data map[string]interface{}
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &data))
	is.Equal(data, testdata)
}

// TestUnwrappedPanics ensures that a helpful panic will occur if
// respond.To is called without the handler being wrapped properly
// with respond.Handler or respond.HandlerFunc.
func TestUnwrappedPanics(t *testing.T) {
	is := is.New(t)
	respond := respond.New()
	w := httptest.NewRecorder()
	r := request()
	is.PanicWith("respond: must wrap with Handler or HandlerFunc", func() {
		respond.To(w, r).With(http.StatusOK, testdata)
	})
}

// TestNow ensures that called Now() will have the completion code
// run immediately - rather than waiting until the function has exited.
// Will also clean up immediately too.
func TestNow(t *testing.T) {
	is := is.New(t)

	respond := respond.New()
	testW := httptest.NewRecorder()
	r := request()

	respond.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respond.To(w, r).With(http.StatusOK, testdata).Now()

		is.Equal(http.StatusOK, testW.Code)
		var data map[string]interface{}
		is.NoErr(json.Unmarshal(testW.Body.Bytes(), &data))
		is.Equal(data, testdata)

	})(testW, r)

}
