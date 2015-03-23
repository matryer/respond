package respond

import (
	"net/http"
	"sync"
)

// TransformFunc functions can transform data before it
// is written.
type TransformFunc func(w http.ResponseWriter, r *http.Request, data interface{}) interface{}

// Transform sets the TransformFunc that will be called before
// data is written.
func Transform(fn TransformFunc) {
	transformLock.Lock()
	transform = fn
	transformLock.Unlock()
}

// DefaultTransformFunc is the default TransformFunc that wraps
// error types.
var DefaultTransformFunc TransformFunc = func(w http.ResponseWriter, r *http.Request, data interface{}) interface{} {
	// transform errors
	if err, ok := data.(error); ok {
		return map[string]interface{}{"error": err.Error()}
	}
	return data
}

var transform = DefaultTransformFunc
var transformLock sync.RWMutex
