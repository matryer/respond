package encoding

import (
	"net/http"
	"strings"
	"sync"

	"github.com/matryer/respond"
)

// Encoders represents a collection of respond.Encoder objects.
type Encoders struct {
	lock     sync.RWMutex
	encoders map[string]respond.Encoder
	Default  respond.Encoder
}

// New makes a new Encoders set.
func New() *Encoders {
	return &Encoders{
		encoders: make(map[string]respond.Encoder),
		Default:  respond.JSON,
	}
}

// EncoderFunc is the Options.Encoder function field that will get the
// appropriate encoder for responding.
func (m *Encoders) EncoderFunc(_ http.ResponseWriter, r *http.Request) respond.Encoder {
	encoder, ok := m.Match(r.Header.Get("Accept"))
	if !ok {
		return m.Default
	}
	return encoder
}

// Add adds an Encoder.
func (m *Encoders) Add(match string, e respond.Encoder) {
	m.lock.Lock()
	m.encoders[match] = e
	m.lock.Unlock()
}

// Del removes an Encoder.
func (m *Encoders) Del(e respond.Encoder) {
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

// Match uses the string to find a matching Encoder.
func (m *Encoders) Match(s string) (respond.Encoder, bool) {
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
