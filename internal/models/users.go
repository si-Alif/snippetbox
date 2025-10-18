package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID int
	Name string
	Email string
	HashedPassword []byte
	Created time.Time
}

// usermodel strut wrapper around a connection pool . This will allow us to perform operations on the database being insidea struct where the connection pool and the handler exists under the same struct
type UserModel struct {
	DB *sql.DB
}


func (um *UserModel) Insert(name string , email string , password string) error {
	return nil
}

func (um *UserModel) Authenticate(email string , password string) (int , error){
	return 0 , nil
}

func (um *UserModel) Exists(id int) (bool , error) {
	return false , nil
}
