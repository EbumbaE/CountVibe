package psql

import (
	_ "github.com/lib/pq"
)

func (d *Postgres) InsertNewUser(id, username, password string) error {
	driverConn := d.driverConn

	dbRequest := `INSERT INTO users (id, username, password) VALUES ($1, $2, $3)`
	_, err := driverConn.Exec(dbRequest, id, username, password)

	return err
}

func (d *Postgres) GetUsername(userID string) (string, error) {
	driverConn := d.driverConn

	dbRequest := `SELECT id, username FROM users WHERE id=$1`
	var username string = ""
	err := driverConn.QueryRow(dbRequest, userID).Scan(&userID, &username)

	if err != nil {
		return "", err
	}

	return username, err
}

func (d *Postgres) GetUserID(username string) (string, error) {
	driverConn := d.driverConn

	dbRequest := `SELECT id, username FROM users WHERE username=$1`
	var userID string = ""
	err := driverConn.QueryRow(dbRequest, username).Scan(&userID, &username)

	if err != nil {
		return "", err
	}

	return userID, err
}

func (d *Postgres) GetLastUserID() (string, error) {
	driverConn := d.driverConn

	dbRequest := `SELECT id FROM users ORDER BY id DESC LIMIT 1`
	var userID string = ""
	err := driverConn.QueryRow(dbRequest).Scan(&userID)

	if err != nil {
		return "", err
	}

	return userID, err

}

func (d *Postgres) GetUserPassword(username string) (string, error) {
	driverConn := d.driverConn

	dbRequest := `SELECT username, password FROM users WHERE username=$1`
	var password string = ""
	err := driverConn.QueryRow(dbRequest, username).Scan(&username, &password)

	if err != nil {
		return "", err
	}

	return password, err

}

func (d *Postgres) CheckUsernameInDB(username string) (bool, error) {
	driverConn := d.driverConn

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

func (d *Postgres) DeleteUser(username string) error { //todo
	driverConn := d.driverConn

	dbRequest := `DELETE FROM users WHERE username=$1`
	_, err := driverConn.Exec(dbRequest, username)

	return err
}

func (d *Postgres) GetAllUsernames() (chan string, error) {
	driverConn := d.driverConn

	dbRequest := `SELECT username FROM users`
	rows, err := driverConn.Query(dbRequest)
	if err != nil {
		return nil, err
	}

	var getUsername string = ""
	returnChan := make(chan string)

	go func() error {
		defer close(returnChan)

		for rows.Next() {
			if err := rows.Scan(&getUsername); err != nil {
				return err
			}
			returnChan <- getUsername
		}

		return nil
	}()

	return returnChan, nil
}
