package scsmiddleware

import (
	"bufio"
	"bytes"
	"net"
	"net/http"
	"time"

	"github.com/Nigel2392/router/v2"
	"github.com/Nigel2392/router/v2/middleware"
	"github.com/Nigel2392/router/v2/request"
	"github.com/alexedwards/scs/v2"
)

type scsRequestSession struct {
	r     *request.Request
	store *scs.SessionManager
}

func (s *scsRequestSession) Get(key string) interface{} {
	return s.store.Get(s.r.Request.Context(), key)
}

func (s *scsRequestSession) Set(key string, value interface{}) {
	s.store.Put(s.r.Request.Context(), key, value)
}

func (s *scsRequestSession) Destroy() error {
	return s.store.Destroy(s.r.Request.Context())
}

func (s *scsRequestSession) Exists(key string) bool {
	return s.store.Exists(s.r.Request.Context(), key)
}

func (s *scsRequestSession) Delete(key string) {
	s.store.Remove(s.r.Request.Context(), key)
}

// Customized version of scs's Middleware function
// This is due to the fact that the original Middleware function
// does not support the router.Handler interface
func SessionMiddleware(store *scs.SessionManager) func(next router.Handler) router.Handler {
	return func(next router.Handler) router.Handler {
		return router.HandleFunc(func(r *request.Request) {
			var token string
			cookie, err := r.Request.Cookie(store.Cookie.Name)
			if err == nil {
				token = cookie.Value
			}
			ctx, err := store.Load(r.Request.Context(), token)
			if err != nil {
				if middleware.DEFAULT_LOGGER != nil {
					middleware.DEFAULT_LOGGER.Error("[%s] Error loading session: %v", r.IP().String(), err)
				}
				store.ErrorFunc(r.Response, r.Request, err)
				return
			}

			bw := &BufferedResponseWriter{ResponseWriter: r.Response}
			sr := r.Request.WithContext(ctx)

			// Store the old response for later
			oldWriter := r.Response
			// Set the buffered writer as the response writer
			r.Response = bw
			// Set the new request with the context
			r.Request = sr
			// Set the session on the request
			r.Session = &scsRequestSession{r: r, store: store}

			next.ServeHTTP(r)

			if sr.MultipartForm != nil {
				sr.MultipartForm.RemoveAll()
			}

			switch store.Status(ctx) {
			case scs.Modified:
				token, expiry, err := store.Commit(ctx)
				if err != nil {
					if middleware.DEFAULT_LOGGER != nil {
						middleware.DEFAULT_LOGGER.Error("[%s] Error committing session: %v", r.IP().String(), err)
					}
					store.ErrorFunc(oldWriter, r.Request, err)
					return
				}
				store.WriteSessionCookie(ctx, oldWriter, token, expiry)
			case scs.Destroyed:
				store.WriteSessionCookie(ctx, oldWriter, "", time.Time{})
			}

			oldWriter.Header().Add("Vary", "Cookie")

			if bw.Code != 0 {
				oldWriter.WriteHeader(bw.Code)
			}
			oldWriter.Write(bw.Buf.Bytes())
		})
	}
}

type BufferedResponseWriter struct {
	http.ResponseWriter
	Buf         bytes.Buffer
	Code        int
	WroteHeader bool
}

func (bw *BufferedResponseWriter) Write(b []byte) (int, error) {
	return bw.Buf.Write(b)
}

func (bw *BufferedResponseWriter) WriteHeader(code int) {
	if !bw.WroteHeader {
		bw.Code = code
		bw.WroteHeader = true
	}
}
func (bw *BufferedResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj := bw.ResponseWriter.(http.Hijacker)
	return hj.Hijack()
}

func (bw *BufferedResponseWriter) Push(target string, opts *http.PushOptions) error {
	if pusher, ok := bw.ResponseWriter.(http.Pusher); ok {
		return pusher.Push(target, opts)
	}
	return http.ErrNotSupported
}
