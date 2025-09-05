package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter , r *http.Request){

	w.Header().Add("Server" , "Go Web Server")

	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	tmpl , err := template.ParseFiles(files...)

	if err != nil {
		app.logger.Error(err.Error() , "method" , r.Method , "uri" , r.URL.RequestURI() )
		http.Error(w , "Internal server error while parsing a static file ..." , http.StatusInternalServerError)
		return
	}

	tmpl_err  := tmpl.ExecuteTemplate(w , "base", nil)

	if tmpl_err != nil {
		app.logger.Error(err.Error() , "method" , r.Method , "uri" , r.URL.RequestURI())
		http.Error(w , "Internal Server Error ..." , http.StatusInternalServerError)
	}

}

// view snippet
func (app	*application) snippetView(w http.ResponseWriter , r *http.Request){

	id , err := strconv.Atoi(r.PathValue("id"))

	if err!= nil || id < 1 {
		http.NotFound(w , r)
		return
	}

	fmt.Fprintf(w , "Display a specific snippet with id  %d" , id)

}

//create snippet
func (app *application) snippetCreate(w http.ResponseWriter , r *http.Request){
	fmt.Fprintln(w , "Display a form to create a new snippet...")
}

func (app *application) snippetCreatePost(w http.ResponseWriter , r *http.Request){

	w.WriteHeader(http.StatusCreated)

	msg := "Creating a snippet in the storage...."
	fmt.Fprintln(w , msg)

}