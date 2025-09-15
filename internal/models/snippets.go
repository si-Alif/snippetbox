package models

import (
	"database/sql"
	"time"
)

// define a struct that holds all the needed data to represent a snippet and store it in the database
type Snippet struct {
	ID int
	Title string
	Content string
	Created time.Time
	Expires time.Time
}

// creAte a wrapper around the sql.DB connection pool
type SnippetModel struct{
	DB *sql.DB
}

// Insert a new snippet into the database
func (m * SnippetModel) Insert(title string , content string , expires int) (int , error){
	return 0 , nil // by default
}

// return a specific snippet based on the ID
func (m *SnippetModel) Get(int) (Snippet , error){
	return  Snippet{} , nil
}

func (m *SnippetModel) Latest() ([]Snippet , error){
	return nil , nil
}

// now we need to inject this SnippetModel wrapper struct into our application in main() function
