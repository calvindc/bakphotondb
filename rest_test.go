package main

import (
	"net/http"
	"testing"
	"time"
)

type mockWriter struct {
}

// Identical to the http.ResponseWriter interface
func (m *mockWriter) Header() http.Header {
	return make(http.Header)
}

// Use EncodeJson to generate the payload, write the headers with http.StatusOK if
// they are not already written, then write the payload.
// The Content-Type header is set to "application/json", unless already specified.
func (m *mockWriter) WriteJson(v interface{}) error {
	return nil
}

// Encode the data structure to JSON, mainly used to wrap ResponseWriter in
// middlewares.
func (m *mockWriter) EncodeJson(v interface{}) ([]byte, error) {
	return nil, nil
}

// Similar to the http.ResponseWriter interface, with additional JSON related
// headers set.
func (m *mockWriter) WriteHeader(int) {

}
func TestLog(t *testing.T) {
	doLog("aa", "1", &mockWriter{}, []byte("aadddd\n"))
	doLog("bb", "1", &mockWriter{}, []byte("aadddd\n"))
	doLog("aa", "3", &mockWriter{}, []byte("aadddd\n"))
	doLog("aa", "1", &mockWriter{}, []byte("aadddd\n"))
	doLog("cc", "1", &mockWriter{}, []byte("aadddd\n"))
	doLog("aa", "1", &mockWriter{}, []byte("aadddd\n"))
	doLog("aa", "2", &mockWriter{}, []byte("aadddd\n"))
	time.Sleep(time.Second)
	doLog("aa", "1", &mockWriter{}, []byte("aadddd\n"))
	doLog("bb", "1", &mockWriter{}, []byte("aadddd\n"))
	doLog("aa", "3", &mockWriter{}, []byte("aadddd\n"))
	doLog("aa", "1", &mockWriter{}, []byte("aadddd\n"))
	doLog("cc", "1", &mockWriter{}, []byte("aadddd\n"))
	doLog("aa", "1", &mockWriter{}, []byte("aadddd\n"))
	doLog("aa", "2", &mockWriter{}, []byte("aadddd\n"))
	time.Sleep(time.Second)
	doLog("aa", "1", &mockWriter{}, []byte("aadddd\n"))
	doLog("bb", "1", &mockWriter{}, []byte("aadddd\n"))
	doLog("aa", "3", &mockWriter{}, []byte("aadddd\n"))
	doLog("aa", "1", &mockWriter{}, []byte("aadddd\n"))
	doLog("cc", "1", &mockWriter{}, []byte("aadddd\n"))
	doLog("aa", "1", &mockWriter{}, []byte("aadddd\n"))
	doLog("aa", "2", &mockWriter{}, []byte("aadddd\n"))
}
