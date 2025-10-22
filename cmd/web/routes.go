package main

import (
	"net/http"
	"github.com/justinas/alice"
)

// Here we centralized all our routes in the routes() methods of our application struct .
// Once called , routes() returns a pointer to a serveMux containing all the routes of our application

// func (app *application) routes() *http.ServeMux {

func (app *application) routes() http.Handler{ // as it will passed in a middleware and they work on return and input type of http.Handler only that's why we returned http.Handler instead of *http.ServeMux which was configured before

	mux := http.NewServeMux()

	//create a file server which will serve contents from the ./ui/static directory
	file_server := http.FileServer(http.Dir("./ui/static/")) // this will a sub-tree path

	// now we need to use the file_server server as the handler for serving file whenever there's a request to a endpoint with prefix /static/
	mux.Handle("GET /static/" , http.StripPrefix("/static" , file_server))
	// a request to /static/favicon.ico --> stripped /static --> result /favicon.ico went to file_server
	// --> file_server looks up at ./ui/static/favicon.ico

	dynamic := alice.New(app.sessionManager.LoadAndSave , noSurf , app.authenticate)

	// ----------------

	// Unauthenticated routes
	mux.Handle("GET /{$}" , dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}" , dynamic.ThenFunc(app.snippetView))

	// All the authentication related routes
	mux.Handle("GET /user/signup" , dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup" , dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login" , dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login" , dynamic.ThenFunc(app.userLoginPost))


	// ----------------
	// authenticated routes

	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /snippet/create" , protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create" , protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /user/logout" , protected.ThenFunc(app.userLogoutPost))

	// ----------------

	// return app.recoverPanic(app.logRequest(commonHeaders(mux))) // returns a http.Handler
	standard := alice.New(app.recoverPanic , app.logRequest , commonHeaders)

	return standard.Then(mux)

}