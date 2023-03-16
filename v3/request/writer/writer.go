package writer

import (
	"bufio"
	"bytes"
	"net"
	"net/http"
)

// BufferedResponseWriter is a buffered response writer.
type ClearableBufferedResponse interface {
	http.ResponseWriter
	Clear()
	Buffer() *bytes.Buffer
	Writer() http.ResponseWriter
	Finalize()
}

// ClearableBufferedResponseWriter is a buffered response writer that can be cleared.
type ClearableBufferedResponseWriter struct {
	http.ResponseWriter
	Buf         *bytes.Buffer
	Code        int
	WroteHeader bool
}

func NewClearable(w http.ResponseWriter) ClearableBufferedResponse {
	return &ClearableBufferedResponseWriter{ResponseWriter: w, Buf: &bytes.Buffer{}}
}

func (bw *ClearableBufferedResponseWriter) Write(b []byte) (int, error) {
	return bw.Buf.Write(b)
}

func (bw *ClearableBufferedResponseWriter) WriteHeader(code int) {
	if !bw.WroteHeader {
		bw.Code = code
		bw.WroteHeader = true
	}
}

func (bw *ClearableBufferedResponseWriter) Writer() http.ResponseWriter {
	return bw.ResponseWriter
}

func (bw *ClearableBufferedResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj := bw.ResponseWriter.(http.Hijacker)
	return hj.Hijack()
}

func (bw *ClearableBufferedResponseWriter) Push(target string, opts *http.PushOptions) error {
	if pusher, ok := bw.ResponseWriter.(http.Pusher); ok {
		return pusher.Push(target, opts)
	}
	return http.ErrNotSupported
}

func (bw *ClearableBufferedResponseWriter) Clear() {
	bw.Buf.Reset()
	bw.Code = 0
	bw.WroteHeader = false
	for k := range bw.Header() {
		bw.Header().Del(k)
	}
}

func (bw *ClearableBufferedResponseWriter) Buffer() *bytes.Buffer {
	return bw.Buf
}

func (bw *ClearableBufferedResponseWriter) Finalize() {
	if bw.Code != 0 {
		bw.ResponseWriter.WriteHeader(bw.Code)
	}
	bw.ResponseWriter.Write(bw.Buf.Bytes())
}
