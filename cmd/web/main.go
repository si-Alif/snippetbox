package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"snippetbox._alif__.net/internal/models"
)

// Define an application struct to hold the dependencies for our application used through out our entire application using struct embedding . This works as central class maintain all the dependencies of our application in a single place (Singleton Pattern applied)
type application struct {
	logger *slog.Logger
	snippets *models.SnippetModel
	template_cache map[string]*template.Template
	formDecoder *form.Decoder
	sessionManager *scs.SessionManager
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
	logger := slog.New(slog.NewTextHandler(os.Stdout  , nil))

	// opens the db connection and stores the connection pool in the db variable
	db , err := openDB(*dsn)

	if err != nil {
		logger.Error(err.Error());
		os.Exit(1)
	}

	// close the db connection
	defer db.Close()

	template_cache , err := newTemplatecache() // cache the template on the server's disk memory

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// using form decoder package for parsing and retrieving form data
	formDecoder := form.NewDecoder()


	// create a new instance of the session manager
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db) // where to store our sessions
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := &application{
		logger : logger,
		snippets : &models.SnippetModel{DB: db}, // create a new instance of the SnippetModel struct with the connection pool as the DB field
		template_cache: template_cache, // added the cached templated in the application dependencies struct
		formDecoder: formDecoder,
		sessionManager: sessionManager,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519 , tls.CurveP256},
	}

	// defining our own http server struct
	srv := &http.Server{
		Addr: *addr,
		Handler: app.routes(),
		ErrorLog: slog.NewLogLogger(logger.Handler() , slog.LevelError),
		TLSConfig: tlsConfig,
		// for what time being the connection will be kept alive before closing for client's idleness . This is important to keep for safe data transfer and nullify security concerns , also we keep the connection alive for some moment cause performing TCP handshakes on every request is expensive and may cause performance issues . So it's good to keep the connection alive for some time before closing it
		IdleTimeout: time.Minute,

		// setting the read and write timeout to 5 and 10 seconds respectively .
		ReadTimeout: 5 * time.Second, // Setting a short ReadTimeout period helps to mitigate the risk from slow-client attacks
		WriteTimeout: 10 * time.Second, 
	}

	// take the HTTP address we got from terminal and show an output message using the custom logger and start the server
	//1️⃣ logger.Info("Starting server on " , "addr" , *addr)
	// 2️⃣ instead of providing the hashmap's key-value pairs like above in a variadic manner , we can use different slog.<data_type>() methods for safer data passing and parsing
	logger.Info("Starting server on :- " , slog.String("addr" , srv.Addr))

	// formerly , all the routes were configured here and the serveMux that was containing all them was passed here
	// err := http.ListenAndServe(*addr , mux)

	// as the route is now abstracted , we now store call the routes() method which returns a pointer to a serveMux containing all the routes
	err = srv.ListenAndServeTLS("./tls/cert.pem" , "./tls/key.pem")

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