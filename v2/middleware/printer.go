package middleware

import (
	"fmt"
	"time"

	"github.com/Nigel2392/router/v2"
	"github.com/Nigel2392/router/v2/request"
)

func Printer(next router.Handler) router.Handler {
	return router.HandleFuncWrapper{F: func(r *request.Request) {
		start := time.Now()
		method := r.Method()
		path := r.Request.URL.Path
		next.ServeHTTP(r)
		fmt.Printf("%s [%s] %s %s\n", start.Format("2006 Jan 02 15:04:05"), method, time.Since(start), path)
	}}
}
