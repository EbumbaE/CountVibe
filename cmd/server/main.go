package main

import (
	"github.com/EbumbaE/CountVibe/internal/certificate"
	"github.com/EbumbaE/CountVibe/internal/config"
	"github.com/EbumbaE/CountVibe/internal/logger"
	"github.com/EbumbaE/CountVibe/internal/server"
	"github.com/EbumbaE/CountVibe/internal/session"
	"github.com/EbumbaE/CountVibe/internal/storage/psql"
)

func main() {

	logger, err := logger.NewLogger("../../internal/logger/l.log")
	if err != nil {
		panic("Create logger " + err.Error())
	}

	conf, err := config.NewConfig()
	if err != nil {
		logger.Error("Parse config: ", err)
		return
	}

	db, err := psql.Init(conf.Database)
	if err != nil {
		logger.Error("Init database: ", err)
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
