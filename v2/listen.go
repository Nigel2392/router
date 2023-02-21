package router

import "net/http"

func (r *Router) Listen(addr string, s ...*http.Server) error {
	var server *http.Server = getServer(addr, s...)
	server.Handler = r
	return server.ListenAndServe()
}

func (r *Router) ListenTLS(addr, certFile, keyFile string, s ...*http.Server) error {
	var server *http.Server = getServer(addr, s...)
	server.Handler = r
	return server.ListenAndServeTLS(certFile, keyFile)
}

func getServer(addr string, s ...*http.Server) *http.Server {
	var server *http.Server
	if len(s) > 0 {
		server = s[0]
	} else {
		server = &http.Server{
			Addr: addr,
		}
	}
	return server
}
