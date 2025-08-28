package main

import (

	"log"
	"net/http"

)


func main(){

	mux := http.NewServeMux()
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
