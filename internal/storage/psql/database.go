package psql

import (
	"database/sql"
)

type Postgres struct {
	driverConn *sql.DB
}

func Init(config Config) (*Postgres, error) {

	connectionInfo := config.Info
	driverConn, err := sql.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}

	return &Postgres{
		driverConn: driverConn,
	}, nil
}

func (d *Postgres) CheckHealth() (bool, error) {
	driverConn := d.driverConn
	err := driverConn.Ping()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (d *Postgres) Close() {
	driverConn := d.driverConn
	driverConn.Close()
}
