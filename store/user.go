package store

import (
	"database/sql"
	"fmt"

	"lineblocs.com/api/model"
)

func getUserFromDB(id int) (*model.User, error) {
	var userId int
	var username string
	var fname string
	var lname string
	var email string
	fmt.Printf("looking up user %d\r\n", id)
	row := db.QueryRow(`SELECT id, username, first_name, last_name, email FROM users WHERE id=?`, id)

	err := row.Scan(&userId, &username, &fname, &lname, &email)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil { //another error
		return nil, err
	}

	return &model.User{Id: userId, Username: username, FirstName: fname, LastName: lname, Email: email}, nil
}
