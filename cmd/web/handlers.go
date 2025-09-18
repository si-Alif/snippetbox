package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"snippetbox._alif__.net/internal/models"

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
		// error handling using custom logger and error handler
		app.logger.Error(err.Error() , "method" , r.Method , "uri" , r.URL.RequestURI() )
		app.serverError(w , r , err)
		return
	}

	tmpl_err  := tmpl.ExecuteTemplate(w , "base", nil)

	if tmpl_err != nil {
		app.logger.Error(tmpl_err.Error() , "method" , r.Method , "uri" , r.URL.RequestURI())
		app.serverError(w, r , tmpl_err)
	}

}

// view snippet
func (app	*application) snippetView(w http.ResponseWriter , r *http.Request){

	id , err := strconv.Atoi(r.PathValue("id"))

	if err!= nil || id < 1 {
		http.NotFound(w , r)
		return
	}

	snippet , err := app.snippets.Get(id)

	if err != nil {
		if(errors.Is(err , models.ErrNoRecord)){
			http.NotFound(w ,r )
		}else {
			app.serverError(w , r , err)
		}

		return
	}

	fmt.Fprintf(w , "%v" , snippet) // write the response data in plain-text http response

}

//create snippet
func (app *application) snippetCreate(w http.ResponseWriter , r *http.Request){
	fmt.Fprintln(w , "Display a form to create a new snippet...")
}

func (app *application) snippetCreatePost(w http.ResponseWriter , r *http.Request){

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id , err := app.snippets.Insert(title , content , expires)

	if err != nil{
		app.serverError(w , r , err)
		return
	}

	// if the snippet is created successfully , redirect the user to view this snippet
	http.Redirect(w , r , fmt.Sprintf("/snippet/view/%d" , id) , http.StatusSeeOther)

}