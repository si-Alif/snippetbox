package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"snippetbox._alif__.net/internal/assert"
)

func TestCommonHeaders(t *testing.T){
	rr := httptest.NewRecorder()

	r , err := http.NewRequest(http.MethodGet , "/" , nil)

	if err != nil {
		t.Fatal(err)
	}

	// create a basic handler that writes OK to the response body so we can use it as the next handler
	next := http.HandlerFunc(func(w http.ResponseWriter , r *http.Request){
		w.Write([]byte("OKey-Dokey pookie "))
	})

	// wrap thr next around the common headers middleware
	commonHeaders(next).ServeHTTP(rr , r)

	res := rr.Result() // retrieve result of the request from the recorder

	expectedValue := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t , res.Header.Get("Content-Security-Policy") , expectedValue)


	expectedValue = "origin-when-cross-origin"
	assert.Equal(t, res.Header.Get("Referrer-Policy"), expectedValue)

	expectedValue = "nosniff"
	assert.Equal(t, res.Header.Get("X-Content-Type-Options"), expectedValue)


	expectedValue = "deny"
	assert.Equal(t, res.Header.Get("X-Frame-Options"), expectedValue)

	expectedValue = "0"
	assert.Equal(t, res.Header.Get("X-XSS-Protection"), expectedValue)

	expectedValue = "Go"
	assert.Equal(t, res.Header.Get("Server"), expectedValue)

	assert.Equal(t , res.StatusCode , http.StatusOK)

	defer res.Body.Close()

	body , err := io.ReadAll(res.Body)

	if err != nil {
		t.Fatal(err)
	}

	body = bytes.TrimSpace(body)

	assert.Equal(t , string(body) , "OKey-Dokey pookie")

}