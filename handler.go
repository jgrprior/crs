package crs

import (
	"io"
	"net/http"
)

// NewHandler returns an HTTP handler function wrapped with auth and validation.
func NewHandler(db Persister, usr string, pass string) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.ContentLength == 0 {
			JSONErrorResponse(w, "Bad request", http.StatusBadRequest)
			return
		}

		defer r.Body.Close()
		entry, err := NewEntry(r.Body)
		if err != nil {
			JSONErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		if valid, errs := entry.Valid(); valid == false {
			JSONErrorResponses(w, errs, http.StatusBadRequest)
			return
		}

		if err := db.Save(&entry); err != nil {
			// XXX: Do proper logging. Don't panic
			// Let the panic handler deal
			panic(err)
		}

		resp := &JSONResponse{Status: http.StatusOK, EntryID: entry.PublicID}
		w.Header().Set("Content-Type", "aplication/json")
		w.WriteHeader(http.StatusOK)
		s, _ := resp.JSONEncode()
		io.WriteString(w, s)
	})

	// TODO: Check content-type
	return WithRecover(WithAuth(usr, pass, WithPost(handler)))
}
