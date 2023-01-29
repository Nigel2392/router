package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Nigel2392/router"
)

func Printer(next router.Handler) router.Handler {
	return router.HandleFuncWrapper{F: func(v router.Vars, w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		method := r.Method
		path := r.URL.Path
		next.ServeHTTP(v, w, r)
		fmt.Printf("%s [%s] %s %s\n", start.Format("2006 Jan 02 15:04:05"), method, time.Since(start), path)
	}}
}
