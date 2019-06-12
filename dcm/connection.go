package dcm

import (
	"errors"

	"github.com/esvm/dcm-middleware/distribution"
	"github.com/esvm/dcm-middleware/distribution/models"
	log "github.com/sirupsen/logrus"
)

type Connection struct {
	requestor  *distribution.Requestor
	aggregator *distribution.Aggregator
}

func (conn *Connection) CreateTopic(name string, properties models.TopicProperties) (*models.Topic, error) {
	request := models.NewCreateTopicRequest(name, properties)
	chRes, err := conn.requestor.Invoke(request)
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	response := <-chRes
	if response.Error() != "" {
		return nil, errors.New(response.Error())
	}

	return response.(*models.CreateTopicResponse).TopicResult, nil
}

func (conn *Connection) Publish(topicID string, message *models.Message) error {
	bufferFull, err := conn.aggregator.AppendMetric(topicID, message)
	if err != nil {
		return err
	}

	if bufferFull {
		metrics := conn.aggregator.Aggregate(topicID)
		request := models.NewPublishRequest(topicID, metrics)
		chRes, err := conn.requestor.Invoke(request)

		if err != nil {
			return err
		}

		response := <-chRes
		return errors.New(response.Error())
	}

	return nil
}

func (conn *Connection) Consume(topicID, indexName string) (chan *models.Metric, error) {
	request := models.NewConsumeRequest(topicID, indexName)
	chRes, err := conn.requestor.Invoke(request)
	if err != nil {
		return nil, err
	}

	ch := make(chan *models.Metric)

	go func() {
		for response := range chRes {
			log.Debug("chegou um consume")
			log.Debug("response: ", response)
			if response.Error() != "" {
				log.Debug(response.Error())
				continue
			}

			metrics := response.(*models.ConsumeResponse).MessageResult
			ch <- metrics
		}
	}()

	return ch, nil
}

func (conn *Connection) ResetIndex(topicName, indexName string) error {
	request := models.NewResetIndexRequest(topicName, indexName)
	chRes, err := conn.requestor.Invoke(request)
	if err != nil {
		return err
	}

	response := <-chRes

	return errors.New(response.Error())
}

func (conn *Connection) Close() error {
	return conn.requestor.Close()
}
