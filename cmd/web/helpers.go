package main

import (
	"fmt"
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


func (app *application) render(w http.ResponseWriter , r *http.Request , status int , page string , data template_data){
	ts , ok := app.template_cache[page] // check if the page is in the cache

	if !ok{
		err := fmt.Errorf("the template %s is not cached" , page)
		app.serverError(w , r , err)
		return
	}

	w.WriteHeader(status) // return response header based on the status code . For instance "200 OK"(success) or "404 Not Found"(when the user tried to do something wrong , server failed then this error page might be rendered with this status code)

	err:= ts.ExecuteTemplate(w , "base" , data)

	if err != nil {
		app.serverError(w , r , err)
	}

}