package dcm

import (
	"encoding/gob"

	"github.com/esvm/dcm-middleware/distribution"
	"github.com/esvm/dcm-middleware/distribution/models"
	log "github.com/sirupsen/logrus"
)

func initLog() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{})
}

func register() {
	gob.Register(&models.CreateTopicRequest{})
	gob.Register(&models.CreateTopicResponse{})
	gob.Register(&models.PublishRequest{})
	gob.Register(&models.PublishResponse{})
	gob.Register(&models.ConsumeRequest{})
	gob.Register(&models.ConsumeResponse{})
}

func Connect(host string, port int, bufferSize int) (*Connection, error) {
	initLog()
	register()
	requestor, err := distribution.NewRequestor(host, port)
	return &Connection{requestor, distribution.NewAggregator(bufferSize)}, err
}
