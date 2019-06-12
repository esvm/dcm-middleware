// package main

// import (
// 	"fmt"

// 	"github.com/esvm/dcm-middleware/dcm"
// 	"github.com/esvm/dcm-middleware/distribution/models"
// 	"github.com/esvm/dcm-middleware/scylla"
// 	log "github.com/sirupsen/logrus"
// )

// func main() {
// 	c, err := dcm.Connect("localhost", 8426, 100)
// 	defer c.Close()
// 	topic, err := c.CreateTopic("Edjan5", models.TopicProperties{IndexName: "xalala", StartFrom: scylla.Begin})
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	consume(c, topic)
// }

// func consume(c *dcm.Connection, topic *models.Topic) {
// 	ch, err := c.Consume(topic.ID, topic.Properties.IndexName)
// 	if err != nil {
// 		log.Debug(err)
// 	}

// 	for metric := range ch {
// 		log.Info(metric)
// 	}
// }
