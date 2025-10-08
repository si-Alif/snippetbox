package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox._alif__.net/internal/models"
	"snippetbox._alif__.net/internal/validator"
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

// as we added the decoder package and for it's reference where to attach certain value from the form in the struct , we added those form name & field key in the end of each one
type snippetCreateFrom struct {
	Title string `form:"title"`
	Content string `form:"content"`
	Expires int `form:"expires"`
	// embedded validator struct which helps us validating and managing form data validation errors
	validator.Validator `form:"-"` // "-" tells decoder to ignore this field
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

	// Manual decoding & without helpers
	//---------------------------------------------------------------
	// // we're using ParseForm() as we're using form markup in the template to get the form data . This parses the form data and then stores them in http.Request instance r as PostForm() map structure .

	// err := r.ParseForm()

	// if err != nil {
	// 	app.clientError(w , http.StatusBadRequest)
	// 	return
	// }

	// // by default , data we retrieve from the map is a string
	// // Our expires field is a number so we need to convert it to an int
	// expires  , err:= strconv.Atoi(r.PostForm.Get("expires"))

	// if err != nil {
	// 	app.clientError(w , http.StatusBadRequest)
	// 	return
	// }

	// form := snippetCreateFrom{
	// 	Title : r.PostForm.Get("title"),
	// 	Content : r.PostForm.Get("content"),
	// 	Expires : expires,
	// }

	// -----------------------------------------------------------------

	var form snippetCreateFrom

	err := app.decodePostForm(r , &form)

	if err != nil {
		app.clientError(w , http.StatusBadRequest)
		return
	}

	
	form.CheckField(validator.NotBlank(form.Title) , "title" , "This field cannot be blank")
	form.CheckField(validator.MaxChars(100 , form.Title) , "title" , "Title must not be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content) , "content" , "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires , 1 ,7 , 365) , "expires" , "This field must be either 1 , 7 or 365")

	if !form.Valid() {
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