package capture

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

var (
	authUser string = os.Getenv("GOCAPTURE_USER")
	authPass string = os.Getenv("GOCAPTURE_PASS")

	// Load our JSON schema once.
	schemaLoader = gojsonschema.NewStringLoader(schema)
)

func init() {
	if authUser == "" || authPass == "" {
		panic("GOCAPTURE_USER and GOCAPTURE_PASS env vars are not set!")
	}
}

// PanicHandler recovers when a handler further down the call chain panics, ensuring
// the error is logged and the HTTP client gets a suitable error response.
func PanicHandler(next http.Handler) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
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
				// TODO: Email alert
				JsonErrorResponse(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}()
		next.ServeHTTP(w, req)
	})
	return handler
}

// PostHandler middleware rejects HTTP requests that do not use the POST method.
func PostHandler(next http.Handler) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Println("Require POST middleware")
		if req.Method != "POST" {
			JsonErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, req)
	})
	return handler
}

// AuthHandler middleware rejects HTTP requests that don't use valid
// basic auth user name and password.
func AuthHandler(next http.Handler) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Println("With auth middleware")

		s := strings.SplitN(req.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 || s[0] != "Basic" {
			JsonErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			JsonErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			JsonErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if pair[0] != authUser || pair[1] != authPass {
			JsonErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, req)
	})
	return handler
}

// ValidJsonHandler middleware rejects HTTP requests that don't contain a valid JSON body
func ValidJsonHandler(next http.Handler) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Println("Validate JSON middleware")

		if req.ContentLength == 0 {
			JsonErrorResponse(w, "Bad request", http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			JsonErrorResponse(w, "Bad request", http.StatusBadRequest)
			return
		}

		// TODO: Look at TeeReader
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(body))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(body))
		bodyCopy, _ := ioutil.ReadAll(rdr1)

		docLoader := gojsonschema.NewStringLoader(string(bodyCopy[:]))
		result, err := gojsonschema.Validate(schemaLoader, docLoader)
		if err != nil {
			JsonErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !result.Valid() {
			errs := make([]string, 0)
			for _, err := range result.Errors() {
				errs = append(errs, err.Description())
				log.Println(err)
			}
			JsonErrorResponses(w, errs, http.StatusBadRequest)
			return
		}

		entry, err := NewEntry(rdr2)
		if err != nil {
			JsonErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(req.Context(), PayloadContextKey, entry)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
	return handler
}
