package main

import (
	"CountVibe/internal/server"
	"CountVibe/internal/config"
	"CountVibe/internal/log"
	"CountVibe/internal/certificate"
)

func main() { 

	logger, err := log.CreateLogger("../../internal/log/l.log")
	if err != nil{
		panic("Create logger error: " + err.Error())
	}
	
	conf := config.CreateConfig()

	certificate.SetupKeyAndCertificate(conf.Certificate)

	serv := server.CreateServer(conf.Server, logger)
	serv.Run(conf.Certificate.Certfile, conf.Certificate.Keyfile)

	

}