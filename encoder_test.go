package respond_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/cheekybits/is"
	"github.com/matryer/respond"
)

func TestJSON(t *testing.T) {
	is := is.New(t)

	w := httptest.NewRecorder()
	r := newTestRequest()

	is.Equal(respond.JSON.ContentType(w, r), "application/json; charset=utf-8")
	respond.JSON.Encode(w, r, testdata)

	var data map[string]interface{}
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &data))
	is.Equal(data, testdata)
}
