package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
)

// to log server side error
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		// debug.Stack() , returns a byte slice which is a stack trace of the execution of code . We need to convert it into a string for readability
		stack_trace = string(debug.Stack())
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", stack_trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

}

// to respond to client side error
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data template_data) {
	ts, ok := app.template_cache[page] // check if the page is in the cache

	if !ok {
		err := fmt.Errorf("the template %s is not cached", page)
		app.serverError(w, r, err)
		return
	}

	// To avoid runtime errors , first of all execute the template in a buffer and based on output decide whether to render that template or not
	buf := new(bytes.Buffer)
	// Execute the template with the data and store the output in the buffer
	err := ts.ExecuteTemplate(buf, "base", data)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status) // return response header based on the status code . For instance "200 OK"(success) or "404 Not Found"(when the user tried to do something wrong , server failed then this error page might be rendered with this status code)

	// if the buffer was correct just write it to the response writer
	buf.WriteTo(w)

}

// this automatically adds these fields in the instance of template_data
func (app *application) newTemplateData(r *http.Request) template_data {
	return template_data{
		Current_year:    time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r), // on instantiation of template_data , this field will be added
	}
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()

	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm) // (<where to store the decoded form data> , <form data map>)
	if err != nil {
		var invalid_decode_error *form.InvalidDecoderError

		if errors.As(err, &invalid_decode_error) {
			panic(err)
		}

		return err
	}

	return nil

}

// Authorization check

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}
