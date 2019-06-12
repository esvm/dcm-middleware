package main

import (
	"fmt"

	"github.com/dwsb/projetomiddleware/dcm"
	"github.com/dwsb/projetomiddleware/distribution/models"
	"github.com/dwsb/projetomiddleware/scylla"
	log "github.com/sirupsen/logrus"
)

func main() {
	c, err := dcm.Connect("localhost", 8426, 100)
	defer c.Close()
	topic, err := c.CreateTopic("Edjan5", models.TopicProperties{IndexName: "xalala", StartFrom: scylla.Begin})
	if err != nil {
		fmt.Println(err)
		return
	}

	consume(c, topic)
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

func consume(c *dcm.Connection, topic *models.Topic) {
	ch, err := c.Consume(topic.ID, topic.Properties.IndexName)
	if err != nil {
		log.Debug(err)
	}

	for metric := range ch {
		log.Info(metric)
	}
}
