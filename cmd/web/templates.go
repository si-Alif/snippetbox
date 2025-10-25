package main

import (
	"html/template"
	"io/fs"
	"path/filepath"

	"time"


	"snippetbox._alif__.net/internal/models"
	"snippetbox._alif__.net/ui"
)


type template_data struct {
	Current_year int //common dynamic data but static in this case or more like cached data . Suppose in the navbar , we show the logo and username of the current user which is dynamic in nature but still we cache ir and then use it in the navbar rather than database or calling on every request
	Snippet models.Snippet
	Snippets []models.Snippet
	Form any // to store form data validation errors
	Flash string // to store flash messages
	IsAuthenticated bool
	CSRFToken string
}

// create a human_date function
func human_date(t time.Time) string{
	if t.IsZero() {
		return ""
	}

	return t.Format("02 Jan 2006 at 15:04")
}

// crate a template.FuncMap with key "human_date" and value human_date function  and store it in a variable named functions .The FuncMap struct is a map that stores all the functions that can be used in templates package and as obvious it's a map structure . So now , we're extending that struct by embedding a new map with key "human_date" and value human_date function and storing it in a variable named functions . Now the functions map is ready to be used in templates cause it ahs all the functionalities of the template.FuncMap as well as the user defined functions values that were embedded to it .
var functions = template.FuncMap{
	"human_date": human_date,
}



func newTemplatecache() (map[string]*template.Template , error){

	//initialize a map where the datatype is a string(name of the file) and the value is a pointer to a template.Template(struct)
	//Template is a specialized Template from "text/template" that produces a safe HTML document fragment.

	cache := map[string]*template.Template{}

	/* ------
	-> we won't cache files anymore by retrieving them from the disk rather we would retrieve them from the embedded file system in the memory

	pages , err := filepath.Glob("./ui/html/pages/*.tmpl.html") // returns a slice of file names matching the pattern "./ui/html/pages/*.tmpl.html"

	*/

	pages , err := fs.Glob(ui.Files , "html/pages/*.tmpl.html")

	if err!= nil{
		return nil , err
	}


	for _ , page := range pages{
		name := filepath.Base(page) // returns the last element of the path means whatever is after the last slash in the path which is indeed the filename

		/*
		--> we won't cache files anymore by retrieving them from the disk rather we would retrieve them from the embedded file system in the memory

		// instantiate a new template set with our custom made template functionalities library
		ts , err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html") // first and foremost parse the base template and store it's struct in ts

		if err != nil{
			return nil , err
		}

		ts , err  = ts.ParseGlob("./ui/html/partials/*.tmpl.html") // parse the partials all the partials available and store them in ts struct

		if err != nil{
			return nil , err
		}

		ts , err = ts.ParseFiles(page) // in the end parse the page and store it's struct in ts . After this ts is now a complete page but cached in the server's disk memory

		if err != nil{
			return nil , err
		}

		cache[name] = ts

		*/

		patterns := []string{
			"html/base.tmpl.html",
			"html/partials/*.tmpl.html",
			page,
		}

		ts , err := template.New(name).Funcs(functions).ParseFS(ui.Files , patterns...)

		if err != nil{
			return nil , err
		}

		cache[name] = ts

	}

	return cache , nil

}


