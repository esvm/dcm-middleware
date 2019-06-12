package broker

import (
	"strings"

	"github.com/dwsb/projetomiddleware/broker/daos"
	"github.com/dwsb/projetomiddleware/distribution/models"
	log "github.com/sirupsen/logrus"
)

type TopicMgmt struct {
	dao daos.TopicDao
}

func NewTopicMgmt() *TopicMgmt {
	return &TopicMgmt{dao: daos.NewTopicDao()}
}

func (t *TopicMgmt) CreateTopic(request *models.CreateTopicRequest) *models.CreateTopicResponse {
	topic := &models.Topic{
		ID:   strings.ToLower(request.Name),
		Name: request.Name,
	}

	topic, err := t.dao.CreateTopic(topic)
	if err != nil {
		log.Debug(err)
		return &models.CreateTopicResponse{
			Err:   err.Error(),
			Alive: false,
		}
	}

	_, err = t.dao.InsertIndex(topic.ID, request.Properties.IndexName, request.Properties.StartFrom)
	if err != nil {
		log.Debug(err)
		return &models.CreateTopicResponse{
			Err:   err.Error(),
			Alive: false,
		}
	}

	topic.Properties = request.Properties

	return &models.CreateTopicResponse{
		TopicResult: topic,
		Err:         "",
		Alive:       false,
	}
}

func (t *TopicMgmt) Publish(request *models.PublishRequest) *models.PublishResponse {
	err := t.dao.InsertMessage(request.TopicID, request.Metrics)
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}

	return &models.PublishResponse{
		Err:   errMsg,
		Alive: false,
	}
}

func (t *TopicMgmt) Consume(request *models.ConsumeRequest, limit int) []*models.ConsumeResponse {
	index, err := t.dao.GetIndex(request.IndexName)
	if err != nil {
		log.Debug("1: ", err)
		return []*models.ConsumeResponse{
			&models.ConsumeResponse{
				Err:   err.Error(),
				Alive: false,
			},
		}
	}

	metrics, err := t.dao.GetMessages(request.TopicID, index, limit)
	if err != nil {
		log.Debug("2: ", err)
		return []*models.ConsumeResponse{
			&models.ConsumeResponse{
				Err:   err.Error(),
				Alive: false,
			},
		}
	}

	_, err = t.dao.IncIndex(request.TopicID, index, len(metrics)-1)
	if err != nil {
		log.Debug("3: ", err)
		return []*models.ConsumeResponse{
			&models.ConsumeResponse{
				Err:   err.Error(),
				Alive: false,
			},
		}
	}

	var response []*models.ConsumeResponse
	for _, metric := range metrics {
		consumeResponse := &models.ConsumeResponse{
			MessageResult: metric,
			Err:           "",
			Alive:         true,
		}

		response = append(response, consumeResponse)
	}

	_, err = t.dao.IncIndex(request.TopicID, index, 1)
	if err != nil {
		log.Debug("1: ", err)
		EOF := &models.ErrEOF{}
		consumeResponse := &models.ConsumeResponse{
			Err:   EOF.Error(),
			Alive: false,
		}

		response = append(response, consumeResponse)
	}

	return response
}
