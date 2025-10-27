package main

import (
	"net/http"
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

	app := newTestApplication(t)

	ts := newTestServer(t , app.routes())

	defer ts.Close()

	code , _ , body := ts.get(t , "/ping")

	assert.Equal(t , code , http.StatusOK)
	assert.Equal(t , body , "OK")

}