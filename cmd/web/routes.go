package main

import "net/http"

// Here we centralized all our routes in the routes() methods of our application struct .
// Once called , routes() returns a pointer to a serveMux containing all the routes of our application

func (app *application) routes() *http.ServeMux {

	mux := http.NewServeMux()

	//create a file server which will serve contents from the ./ui/static directory
	file_server := http.FileServer(http.Dir("./ui/static/")) // this will a sub-tree path

	// now we need to use the file_server server as the handler for serving file whenever there's a request to a endpoint with prefix /static/
	mux.Handle("GET /static/" , http.StripPrefix("/static" , file_server))
	// a request to /static/favicon.ico --> stripped /static --> result /favicon.ico went to file_server
	// --> file_server looks up at ./ui/static/favicon.ico

	mux.HandleFunc("GET /{$}" , app.home)
	mux.HandleFunc("GET /snippet/view/{id}" , app.snippetView)
	mux.HandleFunc("GET /snippet/create" , app.snippetCreate)

	// POST request
	mux.HandleFunc("POST /snippet/create" , app.snippetCreatePost)

	return mux

}