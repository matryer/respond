package respond

import (
	"encoding/json"
	"net/http"
	"strings"

	"sync"
)

// Encoder descirbes an object capable of encoding
// a response.
type Encoder interface {
	// Encode writes a serialization of v to w, optionally using additional
	// information from the http.Request to do so.
	Encode(w http.ResponseWriter, r *http.Request, v interface{}) error
	// ContentType gets a string that will become the Content-Type header
	// when responding through w to the specified http.Request.
	ContentType(w http.ResponseWriter, r *http.Request) string
}

type encodersMap struct {
	lock     sync.RWMutex
	encoders map[string]Encoder
}

func (m *encodersMap) Add(match string, e Encoder) {
	m.lock.Lock()
	m.encoders[match] = e
	m.lock.Unlock()
}

func (m *encodersMap) Del(e Encoder) {
	m.lock.Lock()
	defer m.lock.Unlock()
	found := false
	for k, v := range m.encoders {
		if v == e {
			delete(m.encoders, k)
			found = true
		}
	}
	if !found {
		panic("respond: encoder not found")
	}
}

func (m *encodersMap) Match(s string) (Encoder, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for k, v := range m.encoders {
		if strings.Contains(strings.ToLower(s), strings.ToLower(k)) {
			// match
			return v, true
		}
	}
	return nil, false
}

type jsonEncoder struct{}

var _ Encoder = (*jsonEncoder)(nil)

// JSON is an Encoder for JSON.
var JSON *jsonEncoder

func (*jsonEncoder) Encode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

func (*jsonEncoder) ContentType(w http.ResponseWriter, r *http.Request) string {
	return "application/json"
}
