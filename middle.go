package crs

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
)

// WithRecover recovers when a handler further down the call chain panics, ensuring
// the error is logged and the HTTP client gets a suitable error response.
func WithRecover(next http.Handler) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}
				JSONErrorResponse(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}()
		next.ServeHTTP(w, r)
	})
	return handler
}

// WithPost middleware rejects HTTP requests that do not use the POST method.
func WithPost(next http.Handler) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			JSONErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
	return handler
}

// WithAuth middleware rejects HTTP requests that don't use valid
// basic auth user name and password.
func WithAuth(username string, passwd string, next http.Handler) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 || s[0] != "Basic" {
			JSONErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			JSONErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			JSONErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if pair[0] != username || pair[1] != passwd {
			JSONErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
	return handler
}
