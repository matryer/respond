package encoding_test

import (
	"net/http"
	"testing"

	"github.com/cheekybits/is"
	"github.com/matryer/respond"
	"github.com/matryer/respond/encoding"
)

type testEncoder struct{}

func (testEncoder) Encode(w http.ResponseWriter, r *http.Request, v interface{}) error { return nil }
func (testEncoder) ContentType(w http.ResponseWriter, r *http.Request) string          { return "test/encoder" }

func testRequestWithAccept(accept string) *http.Request {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		panic("bad request: " + err.Error())
	}
	r.Header.Set("Accept", accept)
	return r
}

func TestMatch(t *testing.T) {
	is := is.New(t)

	e1 := &testEncoder{}
	e2 := &testEncoder{}
	e3 := &testEncoder{}

	e := encoding.New()
	e.Default = respond.JSON
	e.Add("json", e1)
	e.Add("XML", e2)

	// Match
	json, ok := e.Match("application/JSON")
	is.True(ok)
	is.Equal(json, e1)
	json = e.EncoderFunc(nil, testRequestWithAccept("application/json"))
	is.Equal(json, e1)
	xml, ok := e.Match("text/xml")
	is.True(ok)
	is.Equal(xml, e2)
	xml = e.EncoderFunc(nil, testRequestWithAccept("text/xml"))
	is.True(ok)
	is.Equal(xml, e2)

	// no responder
	csv, ok := e.Match("text/csv")
	is.False(ok)
	is.Nil(csv)
	csv = e.EncoderFunc(nil, testRequestWithAccept("text/csv"))
	is.Nil(csv)
	is.Equal(e.Default, csv)

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
