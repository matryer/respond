package respond_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/respond"

	"github.com/cheekybits/is"
)

var testdata = map[string]interface{}{"test": true}

type testHandler struct {
	status int
	data   interface{}
}

func (t *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	respond.With(w, r, t.status, t.data)
}

type testStatusHandler struct {
	status int
}

func (t *testStatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	respond.WithStatus(w, r, t.status)
}

func newTestRequest() *http.Request {
	r, err := http.NewRequest("GET", "Something", nil)
	if err != nil {
		panic("bad request: " + err.Error())
	}
	return r
}

func TestWith(t *testing.T) {
	is := is.New(t)

	w := httptest.NewRecorder()
	r := newTestRequest()

	respond.With(w, r, http.StatusOK, testdata)

	is.Equal(http.StatusOK, w.Code)
	var data map[string]interface{}
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &data))
	is.Equal(data, testdata)
	is.Equal(w.HeaderMap.Get("Content-Type"), "application/json; charset=utf-8")
}

func TestWithStatusDefault(t *testing.T) {
	is := is.New(t)

	w := httptest.NewRecorder()
	r := newTestRequest()

	respond.WithStatus(w, r, http.StatusTeapot)

	is.Equal(http.StatusTeapot, w.Code)
	var data map[string]interface{}
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &data))
	is.Equal(data["status"], "I'm a teapot")
	is.Equal(data["code"], http.StatusTeapot)
	is.Equal(w.HeaderMap.Get("Content-Type"), "application/json; charset=utf-8")
}

func TestWithStatusOptions(t *testing.T) {
	is := is.New(t)

	options := &respond.Options{
		StatusData: func(w http.ResponseWriter, r *http.Request, status int) interface{} {
			return map[string]interface{}{"s": status}
		},
	}
	testHandler := &testStatusHandler{
		status: http.StatusTeapot,
	}
	handler := options.Handler(testHandler)

	w := httptest.NewRecorder()
	r := newTestRequest()

	handler.ServeHTTP(w, r)

	is.Equal(http.StatusTeapot, w.Code)
	is.Equal(w.Body.String(), `{"s":418}`+"\n")
	is.Equal(w.HeaderMap.Get("Content-Type"), "application/json; charset=utf-8")
}

func TestBefore(t *testing.T) {
	is := is.New(t)

	options := &respond.Options{}
	var beforecall map[string]interface{}
	options.Before = func(w http.ResponseWriter, r *http.Request, status int, data interface{}) (int, interface{}) {
		beforecall = map[string]interface{}{
			"w": w, "r": r, "status": status, "data": data,
		}
		return status, data
	}

	testHandler := &testHandler{
		status: http.StatusOK, data: testdata,
	}
	handler := options.Handler(testHandler)

	w := httptest.NewRecorder()
	r := newTestRequest()

	handler.ServeHTTP(w, r)

	is.OK(beforecall)
	is.Equal(beforecall["w"], w)
	is.Equal(beforecall["r"], r)
	is.Equal(beforecall["status"], testHandler.status)
	is.Equal(beforecall["data"], testHandler.data)

	is.Equal(http.StatusOK, w.Code)
	var data map[string]interface{}
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &data))
	is.Equal(data, testdata)
	is.Equal(w.HeaderMap.Get("Content-Type"), "application/json; charset=utf-8")

}

func TestAfter(t *testing.T) {
	is := is.New(t)

	options := &respond.Options{}
	var aftercall map[string]interface{}
	options.After = func(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
		aftercall = map[string]interface{}{
			"w": w, "r": r, "status": status, "data": data,
		}
	}

	testHandler := &testHandler{
		status: http.StatusOK, data: testdata,
	}
	handler := options.Handler(testHandler)

	w := httptest.NewRecorder()
	r := newTestRequest()

	handler.ServeHTTP(w, r)

	is.OK(aftercall)
	is.Equal(aftercall["w"], w)
	is.Equal(aftercall["r"], r)
	is.Equal(aftercall["status"], testHandler.status)
	is.Equal(aftercall["data"], testHandler.data)

	is.Equal(http.StatusOK, w.Code)
	var data map[string]interface{}
	is.NoErr(json.Unmarshal(w.Body.Bytes(), &data))
	is.Equal(data, testdata)
	is.Equal(w.HeaderMap.Get("Content-Type"), "application/json; charset=utf-8")

}

type testEncoder struct {
	err error
}

func (e *testEncoder) Encode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	io.WriteString(w, "testEncoder")
	return e.err
}

func (e *testEncoder) ContentType(w http.ResponseWriter, r *http.Request) string {
	return "test/encoder"
}

func TestEncoder(t *testing.T) {
	is := is.New(t)

	options := &respond.Options{}
	var encodercall map[string]interface{}
	options.Encoder = func(w http.ResponseWriter, r *http.Request) respond.Encoder {
		encodercall = map[string]interface{}{
			"w": w, "r": r,
		}
		return &testEncoder{}
	}

	testHandler := &testHandler{
		status: http.StatusOK, data: testdata,
	}
	handler := options.Handler(testHandler)

	w := httptest.NewRecorder()
	r := newTestRequest()

	handler.ServeHTTP(w, r)

	is.OK(encodercall)
	is.Equal(encodercall["w"], w)
	is.Equal(encodercall["r"], r)

	is.Equal(http.StatusOK, w.Code)
	is.Equal(w.Body.String(), "testEncoder")
	is.Equal(w.HeaderMap.Get("Content-Type"), "test/encoder")

}

func TestEncoderOnErr(t *testing.T) {
	is := is.New(t)

	var onErrCall map[string]interface{}
	options := &respond.Options{
		OnErr: func(err error) {
			onErrCall = map[string]interface{}{"err": err}
		},
	}
	encoderErr := errors.New("something went wrong while encoding")
	options.Encoder = func(w http.ResponseWriter, r *http.Request) respond.Encoder {
		return &testEncoder{
			err: encoderErr,
		}
	}

	testHandler := &testHandler{
		status: http.StatusOK, data: testdata,
	}
	handler := options.Handler(testHandler)

	w := httptest.NewRecorder()
	r := newTestRequest()

	handler.ServeHTTP(w, r)

	is.OK(onErrCall)
	is.Equal(onErrCall["err"], encoderErr)

}

func TestOnErrPanic(t *testing.T) {
	is := is.New(t)

	o := respond.Options{}

	err := errors.New("something went wrong")
	is.PanicWith("respond: "+err.Error(), func() { o.OnErrPanic(err) })

}

func TestMultipleWith(t *testing.T) {
	is := is.New(t)

	options := &respond.Options{}
	handler := options.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respond.With(w, r, http.StatusInternalServerError, errors.New("borked"))
		respond.With(w, r, http.StatusOK, nil)
	}))

	w := httptest.NewRecorder()
	r := newTestRequest()

	is.PanicWith("respond: multiple responses", func() {
		handler.ServeHTTP(w, r)
	})

	options = &respond.Options{
		AllowMultiple: true,
	}
	handler = options.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respond.With(w, r, http.StatusInternalServerError, errors.New("borked"))
		respond.With(w, r, http.StatusOK, nil)
	}))

	w = httptest.NewRecorder()
	r = newTestRequest()

	handler.ServeHTTP(w, r)

}
