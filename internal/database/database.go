package database

import(
	"database/sql"
    _ "github.com/lib/pq"
)

type Database struct{
    driverConn *sql.DB
    Info string
}

var database *Database

func Init(config Database)(error){
	
	connectionInfo := config.Info
    driverConn, err := sql.Open("postgres", connectionInfo)
    if err != nil {
        return err
    }   

    database = &Database{
        driverConn: driverConn,
        Info: connectionInfo,
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
       
    dbRequest := `INSERT INTO users (username, password) VALUES ($1, $2)`
    _, err := driverConn.Exec(dbRequest, username, password)

    return err
}

func GetUserPassword(username string)(string, error){
    driverConn := database.driverConn
    
    dbRequest := `SELECT Username, Password FROM users WHERE Username=$1`
    var password string = ""
    err := driverConn.QueryRow(dbRequest, username).Scan(&username, &password)

    if err != nil{
        return "", err
    }
    
    return password, err   

}

func CheckUserInDB(username string)(bool, error){
    driverConn := database.driverConn
    
    dbRequest := `SELECT username FROM users WHERE username=$1`
    rows, err := driverConn.Query(dbRequest, username)
    if err != nil{
        return false, err
    }
    rows.Next()
    
    var getusername string = ""
    if err := rows.Scan(&getusername); err != nil {
        return false, err
    }

    if getusername == username{
        return true, nil 
    }

    return false, nil
}

func DeleteUser(username string)(error){ //todo
    driverConn := database.driverConn
    
    dbRequest := `DELETE FROM users WHERE username=$1`
    _, err := driverConn.Exec(dbRequest, username)

    return err
}