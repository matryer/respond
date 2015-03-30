package respond

import (
	"net/http"
	"sync"
)

// SetHeader specifies a response header.
// Headers set this way will overwrite existing headers.
// See http.Header.Set.
func (with *W) SetHeader(key, value string) *W {
	with.initheaders()
	with.header.Set(key, value)
	return with
}

// AddHeader specifies a response header.
// Headers set this way will append to existing headers.
// See http.Header.Add.
func (with *W) AddHeader(key, value string) *W {
	with.initheaders()
	with.header.Add(key, value)
	return with
}

// DelHeader deletes the specified response header.
// See http.Header.Del.
func (with *W) DelHeader(key string) *W {
	with.initheaders()
	with.header.Del(key)
	return with
}

// initheaders sets up headers for this With copying global
// headers as a starting place.
func (with *W) initheaders() {
	if with.header == nil {
		with.header = make(http.Header)
		if len(headers.headers) > 0 {
			headers.lock.RLock()
			copyheaders(headers.headers, with.header)
			headers.lock.RUnlock()
		}
	}
}

// headers represent global headers, accessible via the
// respond.Headers() function.
var headers = &safeHeader{
	headers: make(http.Header),
}

// Headers gets the http.Header items that will be set on every
// response.
// Use respond.With{}.SetHeader() for response specific headers.
func Headers() interface {
	Add(key, value string)
	Del(key string)
	Get(key string) string
	Set(key, value string)
	Clear()
} {
	return headers
}

const (
	set = true
	add = false
)

func copyheaders(from, to http.Header) {
	if len(from) == 0 {
		return
	}
	for k, vs := range from {
		for _, v := range vs {
			to.Add(k, v)
		}
	}
}

// safeHeader is a concurrent safe http.Header.
type safeHeader struct {
	headers http.Header
	lock    sync.RWMutex
}

func (s *safeHeader) Add(key, value string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.headers.Add(key, value)
}
func (s *safeHeader) Del(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.headers.Del(key)
}
func (s *safeHeader) Get(key string) string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.Get(key)
}
func (s *safeHeader) Set(key, value string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.headers.Set(key, value)
}
func (s *safeHeader) Clear() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.headers = make(http.Header)
}
