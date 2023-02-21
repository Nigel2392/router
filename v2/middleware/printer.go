package middleware

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Nigel2392/router/v2"
	"github.com/Nigel2392/router/v2/request"
)

func PrinterFunc(next router.Handler, out io.Writer) router.Handler {
	return router.ToHandler(func(r *request.Request) {
		start := time.Now()
		method := r.Method()
		path := r.Request.URL.Path
		next.ServeHTTP(r)
		fmt.Fprintf(out, "%s [%s] %s %s\n", start.Format("2006 Jan 02 15:04:05"), method, time.Since(start), path)
	})
}

func Printer(next router.Handler) router.Handler {
	return PrinterFunc(next, os.Stdout)
}
