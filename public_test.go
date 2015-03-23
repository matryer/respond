package respond_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cheekybits/is"
	"github.com/matryer/respond"
)

type privateData struct {
	SecretCode string
	PublicKey  string
}

func (p *privateData) Public() interface{} {
	return map[string]interface{}{"PublicKey": p.PublicKey}
}

func TestPublic(t *testing.T) {
	is := is.New(t)
	w := httptest.NewRecorder()
	r := request()
	obj := &privateData{
		SecretCode: "ABC123",
		PublicKey:  "123456",
	}
	respond.With(
		http.StatusOK,
		obj,
	).To(w, r)
	is.Equal(http.StatusOK, w.Code)
	var data map[string]interface{}
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &data))
	is.Equal(data["PublicKey"], obj.PublicKey)
	is.Nil(data["SecretCode"])
}

type privateData2 struct {
	SecretCode string
	PublicKey  string
}

func (p *privateData2) Public() interface{} {
	return &privateData3{SecretCode: p.SecretCode, PublicKey: p.PublicKey}
}

type privateData3 struct {
	SecretCode string
	PublicKey  string
}

func (p *privateData3) Public() interface{} {
	return map[string]interface{}{"PublicKey": p.PublicKey}
}

func TestNestedPublic(t *testing.T) {
	is := is.New(t)
	w := httptest.NewRecorder()
	r := request()
	obj := &privateData2{
		SecretCode: "ABC123",
		PublicKey:  "123456",
	}
	respond.With(
		http.StatusOK,
		obj,
	).To(w, r)
	is.Equal(http.StatusOK, w.Code)
	var data map[string]interface{}
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &data))
	is.Equal(data["PublicKey"], obj.PublicKey)
	is.Nil(data["SecretCode"])
}
