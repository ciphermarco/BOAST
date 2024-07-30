package api

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ciphermarco/BOAST/log"
	"github.com/go-chi/render"
)

func (env *env) authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		const secretMaxSize = 44

		// 1. Check Authorization header is not empty
		if auth == "" {
			err := errors.New("the Authorization header is missing")
			render.Render(w, r, errUnauthorized(err))
			return
		}

		// 2. Check Authorization header is in the format "<type> <secret>"
		authSplit := strings.Split(auth, " ")
		if len(authSplit) != 2 {
			err := errors.New("wrong authorization format")
			render.Render(w, r, errUnauthorized(err))
			return
		}

		// 3. Check Authorization type is correct (i.e. "Secret") and that
		//    <secret> does not exceed the maximum accepted size in bytes
		authType := authSplit[0]
		b64secret := authSplit[1]
		if authType != "Secret" {
			err := errors.New("unsupported authorization type")
			render.Render(w, r, errUnauthorized(err))
			return
		} else if base64.StdEncoding.DecodedLen(len(b64secret)) > secretMaxSize {
			err := fmt.Errorf("secret is too long; maximum is %d bytes of decoded content", secretMaxSize)
			render.Render(w, r, errUnauthorized(err))
			return
		}

		// 4. Check Authorization is valid base64
		secret, err := base64.StdEncoding.DecodeString(b64secret)
		if err != nil {
			log.Debug("base64 error: %v", err)
			err := errors.New("base64 error")
			render.Render(w, r, errUnauthorized(err))
			return
		}

		// 5. Generate a base32 URL-safe id via SetTest
		id, canary, err := env.strg.SetTest(secret)
		if id == "" || canary == "" || err != nil {
			log.Debug("set test error: %v", err)
			err := fmt.Errorf("could not create test")
			render.Render(w, r, errUnauthorized(err))
			return
		}

		ctx := context.WithValue(r.Context(), idCtxKey, id)
		ctx = context.WithValue(ctx, canaryCtxKey, canary)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
