package main

import (

	"log"
	"net/http"

)


func main(){

	mux := http.NewServeMux()

	//create a file server which will serve contents from the ./ui/static directory
	file_server := http.FileServer(http.Dir("./ui/static/")) // this will a sub-tree path

	// now we need to use the file_server server as the handler for serving file whenever there's a request to a endpoint with prefix /static/
	mux.Handle("GET /static/" , http.StripPrefix("/static" , file_server))
	// a request to /static/favicon.ico --> stripped /static --> result /favicon.ico went to file_server
	// --> file_server looks up at ./ui/static/favicon.ico

	mux.HandleFunc("GET /{$}" , home)
	mux.HandleFunc("GET /snippet/view/{id}" , snippetView)
	mux.HandleFunc("GET /snippet/create" , snippetCreate)

	// POST request
	mux.HandleFunc("POST /snippet/create" , snippetCreatePost)

	log.Print("starting server on http://localhost:8080/")

	err := http.ListenAndServe(":8080" , mux)
	if err!= nil{
		log.Fatal(err)
	}

}
