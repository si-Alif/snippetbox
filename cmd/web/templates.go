package main

import(
	"html/template"
	"path/filepath"
	"snippetbox._alif__.net/internal/models"
)


type template_data struct {
	Snippet models.Snippet
	Snippets []models.Snippet
}


func newTemplatecache() (map[string]*template.Template , error){

	//initialize a map where the datatype is a string(name of the file) and the value is a pointer to a template.Template(struct)
	//Template is a specialized Template from "text/template" that produces a safe HTML document fragment.

	cache := map[string]*template.Template{}

	pages , err := filepath.Glob("./ui/html/pages/*.tmpl.html") // returns a slice of file names matching the pattern "./ui/html/pages/*.tmpl.html"

	if err!= nil{
		return nil , err
	}

	for _ , page := range pages{
		name := filepath.Base(page) // returns the last element of the path means whatever is after the last slash in the path which is indeed the filename

		ts , err := template.ParseFiles("./ui/html/base.tmpl.html") // first and foremost parse the base template and store it's struct in ts

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

	}

	return cache , nil

}