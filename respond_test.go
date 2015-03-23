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

func TestWith(t *testing.T) {
	is := is.New(t)
	w := httptest.NewRecorder()
	r := request()
	respond.With{
		Code: http.StatusOK,
		Data: testdata,
	}.To(w, r)
	is.Equal(http.StatusOK, w.Code)
	var data map[string]interface{}
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &data))
	is.Equal(data, testdata)
}

func TestWithHeader(t *testing.T) {
	is := is.New(t)
	w := httptest.NewRecorder()
	r := request()
	respond.With{
		Code: http.StatusOK,
		Data: testdata,
	}.SetHeader("X-Custom", "Value").To(w, r)
	is.Equal(http.StatusOK, w.Code)
	var data map[string]interface{}
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &data))
	is.Equal(data, testdata)
	is.Equal(w.Header().Get("X-Custom"), "Value")
}

func TestHeadersWithHeader(t *testing.T) {
	is := is.New(t)
	w := httptest.NewRecorder()
	r := request()
	respond.Headers().Add("X-List", "1")
	respond.Headers().Add("X-List", "2")
	respond.Headers().Add("X-List", "3")
	respond.Headers().Set("X-Global", "Value 2")
	respond.Headers().Set("X-Global2", "Value 2")
	respond.Headers().Set("X-Custom", "should be overwritten")
	respond.With{
		Code: http.StatusOK,
		Data: testdata,
	}.
		SetHeader("X-Custom", "overwrite").
		AddHeader("X-List", "4").
		DelHeader("X-Global2").
		To(w, r)
	is.Equal(http.StatusOK, w.Code)
	var data map[string]interface{}
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &data))
	is.Equal(data, testdata)
	is.Equal(w.Header().Get("X-Custom"), "overwrite")
	is.Equal(w.Header().Get("X-Global"), "Value 2")
	is.Equal(w.Header()["X-List"], []string{"1", "2", "3", "4"})
	is.Nil(w.Header()["X-Global2"])
}
