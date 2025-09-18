package models

import (
	"database/sql"
	"errors"
	"time"
)

// define a struct that holds all the needed data to represent a snippet and store it in the database

// ⭐⭐⭐ N.B : As go's philosophy , to export any identifier from a package , it must start with a capital letter . Thus making only the struct in title case will export . In this scenario , other packages can access this struct and declare variable of this struct type but can't directly access the struct fields if the fields are in lower case as they were not exported or exposed in the first place .

// TLDR : struct fields should be exported as well if you want to access and perform operations on them from other packages
type Snippet struct {
	ID int
	Title string
	Content string
	Created time.Time
	Expires time.Time
}

// creAte a wrapper around the sql.DB connection pool
type SnippetModel struct{
	DB *sql.DB // sql.DB embedded
}

// Insert a new snippet into the database . This returns the ID of the newly inserted record

// In MySQL DB the command needs to be performed is :
// OPT-1️⃣ : INSERT INTO snippets (title , content , created , expires) VALUES ('title' , 'content' , NOW() , DATE_ADD(NOW() , INTERVAL 365 DAY))
// The above option is good and correct as well , but it's better to leave the needed data as a placeholder using a ⭐placeholder parameter '?'⭐ and then bind them at the time of execution using the Exec() function . Cause we're not certain what kind of info will be passed from the form during submission . Thus it's better to use the placeholder parameter instead of interpolating the data

// OPT-2️⃣ : INSERT INTO snippets (title , content , created , expires) VALUES (?, ? , NOW() , DATE_ADD(NOW() , INTERVAL 365 DAY))

/*-----------------------------------------------------

⭐ To perform DB operations in go , we have options such as :
- DB.Query() ---> returns multiple row sets of query
- DB.QueryRow() ---> returns a single row of query
- DB.Exec() ---> used to perform insert , update and delete operations

-------------------------------------------------------
*/

// as INSERT is a write operation , we need to use the Exec() function
func (mysql *SnippetModel) Insert(title string , content string , expires int) (int , error){

	// create the sql query to be executed . Placeholder parameters are used for binding title , content , expiry duration respectively
	stmnt := `INSERT INTO snippets (title , content , created , expires) VALUES (? , ? , UTC_TIMESTAMP() , DATE(UTC_TIMESTAMP() + INTERVAL ? DAY))`
	// 💀 remember to check this kind of DB syntax error . Spent 1hr just to debug an error where you forgot to use '+' in this command

	// Place required fields for placeholders in order of the query to be executed . It returns a sql.Result object ,with info about performed operation
	res , err := mysql.DB.Exec(stmnt , title , content , expires)

	if err != nil {
		return 0 , err
	}

	id  , err := res.LastInsertId() // retrieve id from the result object

	if err != nil {
		return 0 , err
	}

	// ID returned by LastInsertId is	a int64 , so we need to convert it to int .
	return int(id) , nil

}

// return a specific snippet based on the ID provided in this method
// To achieve this , we need to perform a SELECT operation on the database which will return a single row of data based on the provided ID
// ✅ command : SELECT id ,title , content , created , expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ? ;

func (mysql *SnippetModel) Get(id int) (Snippet , error){
	stmnt := `SELECT id , title , content , created , expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// the QueryRow() function returns a pointer to sql.Row object which takes the statement and untrusted id parameter as arguments
	row := mysql.DB.QueryRow(stmnt , id)
	// this row object contents all the fields of the selected row

	var snippet Snippet // declare a snippet struct , here we'll store the data of the row we queried . To retrieve required data from the row object , we need to use the Scan() method from the sql.Row object

	// If the struct fields were not exported along with the struct , sql package won't have been able to fill or manipulate them
	err := row.Scan(&snippet.ID , &snippet.Title , &snippet.Content , &snippet.Created , &snippet.Expires)
	// behind the scene Scan() method automatically typecasts SQL data into GOlangs supported data types

	/*-------------------------------------------------
	- CHAR, VARCHAR and TEXT map to string.
	- BOOLEAN maps to bool.
	- INT maps to int;
	- BIGINT maps to int64.
	- DECIMAL and NUMERIC map to float.
	- TIME, DATE and TIMESTAMP map to time.Time.
	-------------------------------------------------*/

	if err != nil{
		if errors.Is(err , sql.ErrNoRows){ // check if there's no such entry in the database
			return Snippet{} , ErrNoRecord
		}else {
			return Snippet{} , err
		}
	}

	return snippet , nil

}

// get latest 10 snippets from the database
func (mysql *SnippetModel) Latest() ([]Snippet , error){
	stmnt :=  `SELECT id , title , content , created , expires FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows , err := mysql.DB.Query(stmnt)

	if err != nil {
		return nil , err
	}

	defer rows.Close() // close the sql.Rows() connection stream before returning from the function

	var snippets []Snippet


	// rows.Next() returns true if there are more rows to read and false otherwise . Inside the loop , we use the Scan() method to populate the snippet struct with the data from the current iteration row . Once there's no more rows to read , the loop ends automatically
	for rows.Next(){
		var s Snippet

		err := rows.Scan(&s.ID , &s.Title  , &s.Content , &s.Created , &s.Expires)

		if err != nil {
			return nil , err
		}

		snippets = append(snippets, s)

	}

	// ⭐ Once the iteration is over make sure to check for the possible error during populating the dataset . Just because the iteration was successful , doesn't mean there was no error while populating the dataset

	if err = rows.Err(); err != nil{
		return nil , err
	}

	return snippets , nil

}

// now we need to inject this SnippetModel wrapper struct into our application in main() function
