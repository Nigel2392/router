package middleware

import (
	"fmt"
	"io"
	"time"

	"github.com/Nigel2392/router/v3"
	"github.com/Nigel2392/router/v3/request"
)

func PrinterFunc(next router.Handler, out io.Writer) router.Handler {
	return router.HandleFunc(func(r *request.Request) {
		printerFunc(next, r, out)
	})
}

func Printer(next router.Handler) router.Handler {
	return router.HandleFunc(func(r *request.Request) {
		printerFunc(next, r, &logger{request: r})
	})
}

func printerFunc(next router.Handler, r *request.Request, out io.Writer) {
	start := time.Now()
	next.ServeHTTP(r)

	fmt.Fprintf(out, "%s %s\n", r.IP(), time.Since(start))
}
