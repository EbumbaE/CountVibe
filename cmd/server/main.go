package main

import (
	"CountVibe/internal/server"
	"CountVibe/internal/config"
	"CountVibe/internal/log"
	"CountVibe/internal/certificate"
	"CountVibe/internal/database"
)

func main() { 

	logger, err := log.CreateLogger("../../internal/log/l.log")
	if err != nil{
		panic("Create logger error: " + err.Error())
	}

	if err := database.Init(); err != nil{	
		logger.Error("Init database error: ", err)
	}
	ok, err := database.CheckHealth()
	if !ok{
		logger.Error("Responce database error", err)
	}
	
	conf := config.CreateConfig()

	certificate.SetupKeyAndCertificate(conf.Certificate)

	serv := server.CreateServer(conf.Server, logger)
	serv.Run(conf.Certificate.Certfile, conf.Certificate.Keyfile)

}