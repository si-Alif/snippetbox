package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"snippetbox._alif__.net/internal/models"

)

func (app *application) home(w http.ResponseWriter , r *http.Request){

	w.Header().Add("Server" , "Go Web Server")

	snippets , err := app.snippets.Latest()

	if err != nil {
		app.serverError(w , r , err)
	}

	data := app.new_template_date(r)
	data.Snippets = snippets

	app.render(w , r , http.StatusOK , "home.tmpl.html" ,data)

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

	data := app.new_template_date(r)
	data.Snippet = snippet

	app.render(w , r , http.StatusOK , "view.tmpl.html"  , data)

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