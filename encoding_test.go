package respond

import (
	"net/http"
	"testing"

	"github.com/cheekybits/is"
)

type testEncoder struct{}

func (testEncoder) Encode(w http.ResponseWriter, r *http.Request, v interface{}) error { return nil }
func (testEncoder) ContentType(w http.ResponseWriter, r *http.Request) string          { return "test/encoder" }

func TestMatch(t *testing.T) {
	is := is.New(t)

	e1 := &testEncoder{}
	e2 := &testEncoder{}
	e3 := &testEncoder{}

	e := &encodersMap{
		encoders: map[string]Encoder{
			"json": e1,
			"XML":  e2,
		},
	}

	// Match
	json, ok := e.Match("application/JSON")
	is.True(ok)
	is.Equal(json, e1)
	xml, ok := e.Match("text/xml")
	is.True(ok)
	is.Equal(xml, e2)
	csv, ok := e.Match("text/csv")
	is.False(ok)
	is.Nil(csv)

	// add
	e.Add("csv", e3)
	csv, ok = e.Match("text/xml")
	is.True(ok)
	is.Equal(csv, e3)

	// remove
	e.Del(e2)
	xml, ok = e.Match("text/xml")
	is.False(ok)
	is.Nil(xml)

}
