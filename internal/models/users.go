package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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
	hashed_password , err := bcrypt.GenerateFromPassword([]byte(password)  ,12)

	if err != nil {
		return err
	}

	stmnt := `INSERT INTO users (name , email , hashed_password , created) VALUES (? , ? , ? , UTC_TIMESTAMP())`

	// pass the byte array of the hashed password to the database as a string
	_, err = um.DB.Exec(stmnt , name , email , string(hashed_password))

	var mySQLerr *mysql.MySQLError

	if err != nil {
		if errors.As(err , &mySQLerr){

			if mySQLerr.Number == 1062 && strings.Contains(mySQLerr.Message , "users_uc_email"){
				return ErrDuplicateEmail
			}

			return  mySQLerr
		}
		return err
	}

	return nil

}


func (um *UserModel) Authenticate(email string , password string) (int , error){

	var (
		id int
		hashed_password []byte
	)

	stmnt := `SELECT id , hashed_password FROM users WHERE email = ?`

	err := um.DB.QueryRow(stmnt , email).Scan(&id , &hashed_password)

	if err!= nil {
		if errors.Is(err , sql.ErrNoRows){
			return 0 , ErrInvalidCredentials
		}else {
			return 0 , err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashed_password , []byte(password))

	if err != nil {
		if errors.Is(err , bcrypt.ErrMismatchedHashAndPassword){
			return  0 , ErrInvalidCredentials
		}else {
			return 0 , err
		}
	}

	return id , nil

}


func (um *UserModel) Exists(id int) (bool , error) {
	var exists bool
	stmnt := `SELECT EXISTS(SELECT true FROM users WHERE id = ?)`

	err := um.DB.QueryRow(stmnt , id).Scan(&exists)

	if err != nil {
		return false , err
	}

	return exists , err

}

