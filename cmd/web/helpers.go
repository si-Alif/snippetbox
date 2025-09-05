package main

import (
	"net/http"
	"runtime/debug"
)

// to log server side error
func (app *application) serverError(w http.ResponseWriter , r *http.Request , err error){
	var (
		method = r.Method
		uri = r.URL.RequestURI()
		// debug.Stack() , returns a byte slice which is a stack trace of the execution of code . We need to convert it into a string for readability
		stack_trace = string(debug.Stack())
	)

	app.logger.Error(err.Error() , "method" , method , "uri" , uri , "trace" , stack_trace)
	http.Error(w , http.StatusText(http.StatusInternalServerError) , http.StatusInternalServerError)

}

// to respond to client side error
func (app *application) clientError(w http.ResponseWriter , status int){
	http.Error(w, http.StatusText(status) , status)
}

