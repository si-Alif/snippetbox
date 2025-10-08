package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"snippetbox._alif__.net/internal/models"
)

func (app *application) home(w http.ResponseWriter , r *http.Request){

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

// define a new struct to store form data validation errors of certain fields
type snippetCreateFrom struct {
	Title string
	Content string
	Expires int
	FieldErrors map[string]string
}

//create snippet
func (app *application) snippetCreate(w http.ResponseWriter , r *http.Request){
	data := app.new_template_date(r)

	data.Form = snippetCreateFrom{
		Expires: 365,
	}

	app.render(w , r , http.StatusOK , "create.tmpl.html" , data)

}

func (app *application) snippetCreatePost(w http.ResponseWriter , r *http.Request){

	// we're using ParseForm() as we're using form markup in the template to get the form data . This parses the form data and then stores them in http.Request instance r as PostForm() map structure .

	err := r.ParseForm()

	if err != nil {
		app.clientError(w , http.StatusBadRequest)
		return
	}

	// by default , data we retrieve from the map is a string
	// Our expires field is a number so we need to convert it to an int
	expires  , err:= strconv.Atoi(r.PostForm.Get("expires"))

	if err != nil {
		app.clientError(w , http.StatusBadRequest)
		return
	}

	form := snippetCreateFrom{
		Title : r.PostForm.Get("title"),
		Content : r.PostForm.Get("content"),
		Expires : expires,
		FieldErrors: map[string]string{},
	}


	if strings.TrimSpace(form.Title) == ""{
		form.FieldErrors["title"] = "This field cannot be blank"

	}else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	if strings.TrimSpace(form.Content) == ""{
		form.FieldErrors["content"] = "This field cannot be blank"

	}

	if expires != 1 && expires != 7 && expires != 365 {
		form.FieldErrors["expires"] = "This field must be 1 or 7 or 365"
	}

	if len(form.FieldErrors) > 0 {
		data := app.new_template_date(r)
		data.Form = form
		app.render(w , r , http.StatusUnprocessableEntity , "create.tmpl.html" , data)
		return
	}

	id , err := app.snippets.Insert(form.Title , form.Content , form.Expires)

	if err != nil {
		app.serverError(w , r , err)
		return
	}

	http.Redirect(w, r , fmt.Sprintf("/snippet/view/%d" , id) , http.StatusSeeOther)

}