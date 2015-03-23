package respond_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cheekybits/is"
	"github.com/matryer/respond"
)

type obj struct {
	Value string `json:"value"`
}

func TestDefaultTransforming(t *testing.T) {

}

func TestTransforming(t *testing.T) {
	is := is.New(t)
	w := httptest.NewRecorder()
	r := request()
	o := &obj{"Hello world"}

	respond.Transform(func(r *http.Request, data interface{}) interface{} {
		switch o := data.(type) {
		case *obj:
			return map[string]interface{}{"object-value": o.Value}
		case error:
			return map[string]interface{}{"error": o.Error()}
		}
		return data
	})

	respond.With(
		http.StatusOK,
		o,
	).To(w, r)
	is.Equal(http.StatusOK, w.Code)
	is.Equal(w.Body.String(), `{"object-value":"Hello world"}`+"\n")
}
