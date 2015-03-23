package respond_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cheekybits/is"
	"github.com/matryer/respond"
)

func TestAfter(t *testing.T) {
	is := is.New(t)
	w := httptest.NewRecorder()
	r := request()

	var afterRes *respond.Response
	var afterReq *http.Request
	var afterStatus int
	var afterData interface{}

	respond.KeepBody(true)
	respond.After(func(w *respond.Response, r *http.Request, status int, data interface{}) {
		afterRes = w
		afterReq = r
		afterStatus = status
		afterData = data
	})
	respond.With(
		http.StatusOK,
		testdata,
	).To(w, r)
	is.Equal(http.StatusOK, w.Code)
	var data map[string]interface{}
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &data))
	is.Equal(data, testdata)

	is.Equal(afterRes.Body().String(), w.Body.String())
	is.Equal(afterData, testdata)
	is.Equal(afterStatus, http.StatusOK)
	is.Equal(afterReq, r)

}
