package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// Define an application struct to hold the dependencies for our application used through out our entire application
type application struct {
	logger *slog.Logger
}


func main(){

	// rather than using hardcoded address to use for the server , let's take the address from the command line flag

	// addr is a string type pointer that stores address to a string storing the value of the flag passed in the command line
	addr := flag.String("addr" , ":4000" , "HTTP network address")

	dsn := flag.String("dsn" , "web:P@33ed_pass@/snippetbox?parseTime=true" , "MySQL data source name");

	flag.Parse()

	// add a custom logger to our application for CLI output instead of using the default logger for the desired outcome
	// ✅logger := slog.New(slog.NewTextHandler(os.Stdout , nil))

	// we can modify this further and add what more info we want in our output
	logger := slog.New(slog.NewJSONHandler(os.Stdout , &slog.HandlerOptions{
		Level: slog.LevelDebug ,
		AddSource: true,
	}))

	db , err := openDB(*dsn)

	if err != nil {
		logger.Error(err.Error());
		os.Exit(1)
	}

	defer db.Close()

	app := &application{
		logger : logger,
	}


	// take the HTTP address we got from terminal and show an output message using the custom logger and start the server
	//1️⃣ logger.Info("Starting server on " , "addr" , *addr)
	// 2️⃣ instead of providing the hashmap's key-value pairs like above in a variadic manner , we can use different slog.<data_type>() methods for safer data passing and parsing
	logger.Info("request received" , slog.String("addr" , ":4000"))

	// formerly , all the routes were configured here and the serveMux that was containing all them was passed here
	// err := http.ListenAndServe(*addr , mux)

	// as the route is now abstracted , we now store call the routes() method which returns a pointer to a serveMux containing all the routes
	err = http.ListenAndServe(*addr , app.routes())

	if err!= nil{
		logger.Error(err.Error())
		// log's Fatal() usually exits the program which is usually abstracted from the user . But as we're using our custom logger , we need to terminate our application manually by using the os.Exit(1) , here the 1 is a flag of saying the code was terminated with an error
		os.Exit(1)
	}

}


func openDB(dsn string)(*sql.DB , error){
	db , err := sql.Open("mysql" , dsn)

	if err != nil {
		return nil , err
	}

	err = db.Ping() // to verify if the connection stream is alive and connect if not
	if err != nil {
		return nil , err
	}

	return db , nil

}