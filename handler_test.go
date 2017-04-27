package crs

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// mockDB implements Save and Close methods of the crs.Persister interface.
// mockDB never returns an error.
type mockDB struct{}

func (m *mockDB) Save(e *Entry) error { return nil }
func (m *mockDB) Close()              {}

// mockErrorDB implements Save and Close methods of the crs.Persister interface.
// mockErrorDB always returns an error from Save.
type mockErrorDB struct{}

func (m *mockErrorDB) Save(e *Entry) error { return errors.New("database offline") }
func (m *mockErrorDB) Close()              {}

func TestCRSHandler(t *testing.T) {
	var handler http.Handler

	urlStr := "/campaign"
	authUser := "user"
	authPass := "pass"

	var tests = []struct {
		method         string
		db             Persister
		usr            string
		psw            string
		body           string
		wantStatusCode int
	}{
		{"POST", &mockDB{}, authUser, authPass, "", http.StatusBadRequest},         // Empty body
		{"POST", &mockDB{}, "foo", authPass, "", http.StatusUnauthorized},          // Wrong user
		{"POST", &mockDB{}, authUser, "bar", "", http.StatusUnauthorized},          // Wrong pasword
		{"GET", &mockDB{}, authUser, authPass, "", http.StatusMethodNotAllowed},    // Wong method
		{"PUT", &mockDB{}, authUser, authPass, "", http.StatusMethodNotAllowed},    // Wong method
		{"PATCH", &mockDB{}, authUser, authPass, "", http.StatusMethodNotAllowed},  // Wong method
		{"DELETE", &mockDB{}, authUser, authPass, "", http.StatusMethodNotAllowed}, // Wong method

		{"POST", &mockErrorDB{}, authUser, authPass, `{
			"campaignVersion": "0.0.0",
			"campaignName": "Unit Test Campaign",
			"entrant": {
				"title": "Mr",
				"firstName": "John",
				"lastName": "Smith",
				"emailAddress": "foo@gmail.com"
			},
			"form": []
		}`, http.StatusInternalServerError}, // Database error

		{"POST", &mockDB{}, authUser, authPass, `{
			"campaignVersion": "0.0.0",
			"campaignName": "Unit Test Campaign",
			"entrant": {
				"title": "Mr",
				"firstName": "John",
				"lastName": "Smith",
				"emailAddress": "foo@gmail.com"
			},
			"form": []
		}`, http.StatusOK}, // Good request

		{"POST", &mockDB{}, authUser, authPass, `{
			"campaignVersion": "0.0.0",
			"campaignName": 100,
			"entrant": {
				"title": "Mr",
				"firstName": "John",
				"lastName": "Smith",
				"emailAddress": "foo@gmail.com"
			},
			"form": []
		}`, http.StatusBadRequest}, // Bad request. Campaign name is an integer.

	}

	for _, tt := range tests {
		handler = NewHandler(tt.db, authUser, authPass)
		req, err := http.NewRequest(tt.method, urlStr, strings.NewReader(tt.body))
		if err != nil {
			t.Fatal(err)
		}
		req.SetBasicAuth(tt.usr, tt.psw)

		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		resp := w.Result()

		defer resp.Body.Close()
		jsonResp, err := decodeResponse(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != tt.wantStatusCode {
			t.Errorf(
				"got status code %v and message '%s', expected %v",
				resp.StatusCode, strings.Join(jsonResp.Messages, ", "), tt.wantStatusCode)
		}
	}
}

// decodeResponse unmarshals our ResponseRecorder body into a JSONError
func decodeResponse(body io.Reader) (JSONError, error) {
	decoder := json.NewDecoder(body)
	var resp JSONError
	err := decoder.Decode(&resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
