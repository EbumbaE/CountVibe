package main

import (
	"CountVibe/internal/server"
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

	if err := database.Init(); err != nil{	
		logger.Error("Init database ", err)
	}
	ok, err := database.CheckHealth()
	if !ok{
		logger.Error("Responce database ", err)
	}
	
	conf := config.NewConfig()

	if err := certificate.SetupKeyAndCertificate(conf.Certificate); err != nil{
		logger.Error("Setup certificate ", err)
	}

	serv := server.NewServer(conf.Server, logger)
	serv.Run(conf.Certificate.Certfile, conf.Certificate.Keyfile)

}