package database

import(
	"database/sql"
    _ "github.com/lib/pq"
)

type Database struct{
    driverConn *sql.DB
    info string

}

var database *Database

func Init()(error){
	
	connectionInfo := "host=localhost port=5432 user=postgres password=1234 dbname=abobus sslmode=disable"
    driverConn, err := sql.Open("postgres", connectionInfo)
    if err != nil {
        return err
    }   

    database = &Database{
        driverConn: driverConn,
        info: connectionInfo,
    }

    return err
}

func CheckHealth()(bool, error){
    driverConn := database.driverConn
    err := driverConn.Ping()
    if err != nil {
        return false, err
    }   
    return true, nil
}

func Close(){
    driverConn := database.driverConn
    driverConn.Close()
}

func InsertNewUser(username, password string)(error){
    driverConn := database.driverConn
       
    dbRequest := `INSERT INTO users (Username, Password) VALUES ($1, $2)`
    _, err := driverConn.Exec(dbRequest, username, password)

    return err
}

func GetUserPassword(username string)(string, error){
    driverConn := database.driverConn
    
    dbRequest := `SELECT Username, Password FROM users WHERE Username = $1`
    password := ""
    err := driverConn.QueryRow(dbRequest, username).Scan(&username, &password)

    if err != nil{
        return "", err
    }
    
    return password, err   

}

func CheckUserInDB(username string)(bool, error){
    driverConn := database.driverConn
    
    dbRequest := `SELECT Username FROM users WHERE Username = $1`
    err := driverConn.QueryRow(dbRequest, username).Scan(&username)

    if err != nil{
        return false, err
    }

    if username != ""{
        return true, err
    }

    return false, err
}

func DeleteUser()(error){ //todo
    dbRequest := `DELETE FROM users WHERE ID=$1`
    driverConn := database.driverConn
    _, err := driverConn.Exec(dbRequest, 1)
    return err
}