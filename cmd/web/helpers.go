package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
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

	// To avoid runtime errors , first of all execute the template in a buffer and based on output decide whether to render that template or not
	buf := new(bytes.Buffer)
	// Execute the template with the data and store the output in the buffer
	err := ts.ExecuteTemplate(buf , "base" , data)

	if err != nil {
		app.serverError(w , r , err)
		return
	}

	w.WriteHeader(status) // return response header based on the status code . For instance "200 OK"(success) or "404 Not Found"(when the user tried to do something wrong , server failed then this error page might be rendered with this status code)

	// if the buffer was correct just write it to the response writer
	buf.WriteTo(w)

}

func (app *application) new_template_date(r *http.Request) template_data{
	return template_data{
		Current_year: time.Now().Year(),
	}
}

