package main

import (
	"encoding/gob"
	"fmt"

	"github.com/dwsb/projetomiddleware/dcm"
	"github.com/dwsb/projetomiddleware/distribution/models"
	"github.com/dwsb/projetomiddleware/scylla"
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

func main() {
	initLog()
	register()

	c, err := dcm.Connect("localhost", 8426)
	defer c.Close()
	topic, err := c.CreateTopic("Edjan5", models.TopicProperties{IndexName: "xalala", StartFrom: scylla.Begin})
	if err != nil {
		fmt.Println(err)
		return
	}

	publish(c, topic)
}

func publish(c *dcm.Connection, topic *models.Topic) {
	message := &models.Message{
		TopicID: topic.ID,
		Payload: 1.5,
	}

	err := c.Publish(topic.ID, message)
	if err != nil {
		log.Debug(err)
	}

	message = &models.Message{
		TopicID: topic.ID,
		Payload: 4.7,
	}

	err = c.Publish(topic.ID, message)
	if err != nil {
		log.Debug(err)
	}
}
