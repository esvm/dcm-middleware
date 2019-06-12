package main

import (
	"encoding/gob"
	"os"

	"github.com/dwsb/projetomiddleware/broker"
	"github.com/dwsb/projetomiddleware/distribution/models"
	log "github.com/sirupsen/logrus"
)

func initLog() {
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	lvl := os.Getenv("LOG_LEVEL")
	var level log.Level
	switch lvl {
	case "PANIC":
		level = log.PanicLevel
	case "FATAL":
		level = log.FatalLevel
	case "ERROR":
		level = log.ErrorLevel
	case "WARN":
		level = log.WarnLevel
	case "INFO":
		level = log.InfoLevel
	case "DEBUG":
		level = log.DebugLevel
	case "TRACE":
		level = log.TraceLevel
	default:
		level = log.InfoLevel
	}

	log.SetLevel(level)
}

func register() {
	gob.Register(&models.CreateTopicRequest{})
	gob.Register(&models.CreateTopicResponse{})
	gob.Register(&models.PublishRequest{})
	gob.Register(&models.PublishResponse{})
	gob.Register(&models.ConsumeRequest{})
	gob.Register(&models.ConsumeResponse{})
}

func main() {
	initLog()
	register()

	topicMgmt := broker.NewTopicMgmt()
	sessionMgmt := broker.NewSessionMgmt()
	exchange := broker.NewExchange(topicMgmt, sessionMgmt)
	exchange.Exchange()
}
