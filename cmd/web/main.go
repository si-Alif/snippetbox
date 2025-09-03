package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)


func main(){

	// rather than using hardcoded address to use for the server , let's take the address from the command line flag

	// addr is a string type pointer that stores address to a string storing the value of the flag passed in the command line
	addr := flag.String("addr" , ":4000" , "HTTP network address")

	flag.Parse()

	// add a custom logger to our application for CLI output instead of using the default logger for the desired outcome
	// ✅logger := slog.New(slog.NewTextHandler(os.Stdout , nil))

	// we can modify this further and add what more info we want in our output
	logger := slog.New(slog.NewJSONHandler(os.Stdout , &slog.HandlerOptions{
		Level: slog.LevelDebug ,
		AddSource: true,
	}))

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

	// take the HTTP address we got from terminal and show an output message using the custom logger and start the server
	//1️⃣ logger.Info("Starting server on " , "addr" , *addr)
	// 2️⃣ instead of providing the hashmap's key-value pairs like above in a variadic manner , we can use different slog.<data_type>() methods for safer data passing and parsing
	logger.Info("request received" , slog.String("addr" , ":4000"))

	err := http.ListenAndServe(*addr , mux)
	if err!= nil{
		logger.Error(err.Error())
		// log's Fatal() usually exits the program which is usually abstracted from the user . But as we're using our custom logger , we need to terminate our application manually by using the os.Exit(1) , here the 1 is a flag of saying the code was terminated with an error
		os.Exit(1)
	}

}
