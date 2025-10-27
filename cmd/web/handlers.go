package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox._alif__.net/internal/models"
	"snippetbox._alif__.net/internal/validator"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, r, err)
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, r, http.StatusOK, "home.tmpl.html", data)

}

// view snippet
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippets.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	// retrieve request associated key-value pairs from r.Context() using sessionManager's PopString() middleware function which retrieves the valye of the key "flash" from the session and then deletes the key-value pair from the session

	data := app.newTemplateData(r)
	data.Snippet = snippet

	// ---------------------------------------------------
	// flash := app.sessionManager.PopString(r.Context(), "flash")
	// data.Flash = flash --> Flash been added via newTemplateData()
	// ---------------------------------------------------


	app.render(w, r, http.StatusOK, "view.tmpl.html", data)

}

// define a new struct to store form data validation errors of certain fields

// as we added the decoder package and for it's reference where to attach certain value from the form in the struct , we added those form name & field key in the end of each one
type snippetCreateFrom struct {
	Title   string `form:"title"`
	Content string `form:"content"`
	Expires int    `form:"expires"`
	// embedded validator struct which helps us validating and managing form data validation errors
	validator.Validator `form:"-"` // "-" tells decoder to ignore this field
}

// create snippet
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = snippetCreateFrom{
		Expires: 365,
	}

	app.render(w, r, http.StatusOK, "create.tmpl.html", data)

}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

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

	err := app.decodePostForm(r, &form)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(100, form.Title), "title", "Title must not be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must be either 1 , 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// r.Context() is a method that returns the Context associated with the Request. A Context is an immutable object that contains values for the current request and methods to store values that can be accessed by the handlers for the current request.

	// Here, we use Put() method of the session manager to store a key-value pair in the session. The key is "flash" and the value is "Snippet Created Successfully". The value will be available in the next request and can be accessed using the session manager's Get() method.

	app.sessionManager.Put(r.Context(), "flash", "Snippet Created Successfully")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

}



// ALL about authentication system goes down here

type userSignupForm struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Password string `form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email string `form:"email"`
	Password string `form:"password"`
	validator.Validator `form:"-"`
}


func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}

	app.render(w , r , http.StatusOK , "signup.tmpl.html" , data )

}

func (app *application) userSignupPost(w http.ResponseWriter , r *http.Request){
	var form userSignupForm

	err := app.decodePostForm(r , &form)

	if err!= nil {
		app.clientError(w , http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name) , "name" , "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email) , "email" , "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email , validator.EmailRX) , "email" , "This field must be a valid email address")

	form.CheckField(validator.NotBlank(form.Password) , "password" , "This field cannot be blank")
	form.CheckField(validator.MinChars(8 , form.Password) , "password" , "This field must be at least 8 characters long")

	if !form.Valid(){
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w , r , http.StatusUnprocessableEntity , "signup.tmpl.html" , data)
		return
	}

	err = app.users.Insert(form.Name , form.Email , form.Password)

	if err != nil {
		if errors.Is(err , models.ErrDuplicateEmail){
			form.AddFieldError("email" , "Email address already in use")
			data := app.newTemplateData(r)

			data.Form = form

			app.render(w , r , http.StatusUnprocessableEntity , "signup.tmpl.html" , data)
		}else {
			app.serverError(w , r , err)
		}
		return
	}

	app.sessionManager.Put(r.Context() , "flash" , "You've Successfully been signed up in snippetbox")

	http.Redirect(w , r , "/user/login" , http.StatusSeeOther)

}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w , r , http.StatusOK , "login.tmpl.html" , data)

}

func (app *application) userLoginPost(w http.ResponseWriter , r *http.Request){
	var form userLoginForm

	err := app.decodePostForm(r , &form)

	if err != nil{
		app.clientError(w , http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email) , "email" , "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email , validator.EmailRX) , "email" , "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password) , "password" , "This field cannot be blank")

	if !form.Valid(){
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w , r , http.StatusUnprocessableEntity , "login.tmpl.html" , data)
		return
	}

	id , err := app.users.Authenticate(form.Email , form.Password)

	if err != nil {
		if errors.Is(err , models.ErrInvalidCredentials){
			form.AddNonFieldError("Invalid credentials . Enter correct email and Password again")

			data := app.newTemplateData(r);
			data.Form = form
			app.render(w , r , http.StatusUnprocessableEntity , "login.tmpl.html" , data)
		}else {
			app.serverError(w , r , err)
		}

		return
	}

	err = app.sessionManager.RenewToken(r.Context())

	if err != nil {
		app.serverError(w , r , err)
		return
	}

	app.sessionManager.Put(r.Context() , "authenticatedUserID" , id)

	http.Redirect(w , r , "/snippet/create" , http.StatusSeeOther)

}

func (app *application) userLogoutPost(w http.ResponseWriter , r *http.Request){
	err := app.sessionManager.RenewToken(r.Context())

	if err != nil {
		app.serverError(w , r , err)
		return
	}

	// remove the authenticatedUserID from the session to sign the user out
	app.sessionManager.Remove(r.Context() , "authenticatedUserID")

	app.sessionManager.Put(r.Context() , "flash" , "You've been logged out successfully")

	http.Redirect(w , r , "/" , http.StatusSeeOther)

}


func ping(w http.ResponseWriter , r *http.Request) {
	w.Write([]byte("OK"))
}


