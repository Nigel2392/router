package csrf

import (
	"context"
	"fmt"

	"github.com/Nigel2392/router/v2/request"
)

type contextKey string

const usedKey = contextKey(CSRF_TOKEN_COOKIE_NAME)

type contextValue struct {
	token string
}

func contextSaveToken(r *request.Request, token string) {
	ctx := r.Request.Context()
	ctx.Value(usedKey).(*contextValue).token = token
}

func contextGetToken(r *request.Request) (string, error) {
	val := r.Request.Context().Value(usedKey)
	if val == nil {
		return "", fmt.Errorf("no value exists in the context for key %q", CSRF_TOKEN_COOKIE_NAME)
	}
	return val.(*contextValue).token, nil
}

func defaultContext(r *request.Request) {
	r.Request = r.Request.WithContext(context.WithValue(r.Request.Context(), usedKey, &contextValue{}))
}
