package main

import (
	"CountVibe/internal/certificate"
	"CountVibe/internal/config"
	"CountVibe/internal/log"
	"CountVibe/internal/server"
	"CountVibe/internal/session"
	"CountVibe/internal/storage/psql"
)

func main() {

	logger, err := log.NewLogger("../../internal/log/l.log")
	if err != nil {
		panic("Create logger " + err.Error())
	}

	conf := config.NewConfig()

	db, err := psql.Init(conf.Database)
	if err != nil {
		logger.Error("Init database ", err)
	}
	ok, err := db.CheckHealth()
	if !ok || err != nil {
		logger.Error("Database check health", err)
	}

	if err := certificate.SetupKeyAndCertificate(conf.Certificate); err != nil {
		logger.Error("Setup certificate ", err)
	}

	s := session.NewSession(conf.Session, conf.Pages, db, logger)
	s.Run()

	serv := server.NewServer(conf.Server, conf.Pages, logger)
	serv.Run(conf.Certificate.CertPath, conf.Certificate.KeyPath)

}
