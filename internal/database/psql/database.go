package psql

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Database struct {
	driverConn *sql.DB
	Info       string
}

var database *Database

func Init(config Database) error {

	connectionInfo := config.Info
	driverConn, err := sql.Open("postgres", connectionInfo)
	if err != nil {
		return err
	}

	database = &Database{
		driverConn: driverConn,
		Info:       connectionInfo,
	}

	return err
}

func CheckHealth() (bool, error) {
	driverConn := database.driverConn
	err := driverConn.Ping()
	if err != nil {
		return false, err
	}
	return true, nil
}

func Close() {
	driverConn := database.driverConn
	driverConn.Close()
}

func InsertNewUser(id, username, password string) error {
	driverConn := database.driverConn

	dbRequest := `INSERT INTO users (id, username, password) VALUES ($1, $2, $3)`
	_, err := driverConn.Exec(dbRequest, id, username, password)

	return err
}

func GetUsername(userID string) (string, error) {
	driverConn := database.driverConn

	dbRequest := `SELECT id, username FROM users WHERE id=$1`
	var username string = ""
	err := driverConn.QueryRow(dbRequest, userID).Scan(&userID, &username)

	if err != nil {
		return "", err
	}

	return username, err
}

func GetUserID(username string) (string, error) {
	driverConn := database.driverConn

	dbRequest := `SELECT id, username FROM users WHERE username=$1`
	var userID string = ""
	err := driverConn.QueryRow(dbRequest, username).Scan(&userID, &username)

	if err != nil {
		return "", err
	}

	return userID, err
}

func GetUserPassword(username string) (string, error) {
	driverConn := database.driverConn

	dbRequest := `SELECT username, password FROM users WHERE username=$1`
	var password string = ""
	err := driverConn.QueryRow(dbRequest, username).Scan(&username, &password)

	if err != nil {
		return "", err
	}

	return password, err

}

func CheckUsernameInDB(username string) (bool, error) {
	driverConn := database.driverConn

	dbRequest := `SELECT username FROM users WHERE username=$1`
	rows, err := driverConn.Query(dbRequest, username)
	if err != nil {
		return false, err
	}
	rows.Next()

	var getUsername string = ""
	if err := rows.Scan(&getUsername); err != nil {
		return false, err
	}

	if getUsername == username {
		return true, nil
	}

	return false, nil
}

func DeleteUser(username string) error { //todo
	driverConn := database.driverConn

	dbRequest := `DELETE FROM users WHERE username=$1`
	_, err := driverConn.Exec(dbRequest, username)

	return err
}

func GetAllUsernames() ([]string, error) { //so bad
	driverConn := database.driverConn

	dbRequest := `SELECT username FROM users`
	rows, err := driverConn.Query(dbRequest)
	if err != nil {
		return nil, err
	}

	var getUsername string = ""
	var usernames []string

	for rows.Next() {
		if err := rows.Scan(&getUsername); err != nil {
			return nil, err
		}
		usernames = append(usernames, getUsername)
	}

	return usernames, nil

}
