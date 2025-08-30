package main

import (
	"flag"
	"log"
	"net/http"
)


func main(){

	// rather than using hardcoded address to use for the server , let's take the address from the command line flag

	// addr is a string type pointer that stores address to a string storing the value of the flag passed in the command line
	addr := flag.String("addr" , ":4000" , "HTTP network address")

	flag.Parse()

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

	log.Printf("starting server on , http://localhost%s/" , *addr)

	err := http.ListenAndServe(*addr , mux)
	if err!= nil{
		log.Fatal(err)
	}

}
