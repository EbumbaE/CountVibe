package main

import (
	"CountVibe/internal/server"
	"CountVibe/internal/middleware"
	"CountVibe/internal/config"
	"CountVibe/internal/log"
	"CountVibe/internal/certificate"
	"CountVibe/internal/database"
)

func main() { 


	logger, err := log.NewLogger("../../internal/log/l.log")
	if err != nil{
		panic("Create logger " + err.Error())
	}

	conf := config.NewConfig()

	if err := database.Init(conf.Database); err != nil{	
		logger.Error("Init database ", err)
	}
	ok, err := database.CheckHealth()
	if !ok{
		logger.Error("Responce database ", err)
	}

	if err := certificate.SetupKeyAndCertificate(conf.Certificate); err != nil{
		logger.Error("Setup certificate ", err)
	}

	mw := middleware.NewMiddleware(conf.Middleware, conf.Pages)
	mw.Run()

	serv := server.NewServer(conf.Server, conf.Pages, logger)
	serv.Run(conf.Certificate.Certfile, conf.Certificate.Keyfile)

}