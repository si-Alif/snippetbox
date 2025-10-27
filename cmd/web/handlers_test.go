package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"snippetbox._alif__.net/internal/assert"
)

/* ----- Unit Tests


func TestPing(t *testing.T){
	rr := httptest.NewRecorder() // new recorder to record HTTP response instead of using htt.ResponseWriter to get the output out of the handler

	r , err := http.NewRequest(http.MethodGet, "/" , nil) // make a response to "/" with nil body

	if err != nil {
		t.Fatal(err)
	}

	ping(rr , r) // call the ping handler with the recorder as the writer where the response will be stored and the request

	rs := rr.Result() // retrieve the response

	assert.Equal(t , rs.StatusCode , http.StatusOK) // check if the status code is 200

	defer rs.Body.Close() // stop reading the response body once testing id done

	body , err := io.ReadAll(rs.Body) // read the response body

	if err != nil {
		t.Fatal(err)
	}

	body = bytes.TrimSpace(body) // remove trailing spaces if any from the body

	assert.Equal(t , string(body) , "OK")

}

------ */

func TestPing(t *testing.T){
	app := &application{
		logger: slog.New(slog.DiscardHandler),
	}

	ts := httptest.NewTLSServer(app.routes())

	defer ts.Close()

	res , err := ts.Client().Get(ts.URL + "/ping")

	if err !=nil {
		t.Fatal(err)
	}

	assert.Equal(t , res.StatusCode , http.StatusOK)
	defer res.Body.Close()

	body , err := io.ReadAll(res.Body)

	if err != nil {
		t.Fatal(err)
	}

	body = bytes.TrimSpace(body)
	
	assert.Equal(t , string(body) , "OK")


}