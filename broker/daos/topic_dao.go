package daos

import (
	"github.com/esvm/dcm-middleware/distribution/models"
	"github.com/esvm/dcm-middleware/scylla"
)

type TopicDao interface {
	CreateTopic(*models.Topic) (*models.Topic, error)
	InsertMessage(string, *models.Metric) error
	InsertIndex(string, string, string) (*models.Index, error)
	IncIndex(string, *models.Index, int) (*models.Index, error)
	GetIndex(string) (*models.Index, error)
	GetMessage(string, *models.Index) (*models.Metric, error)
	GetMessages(string, *models.Index, int) ([]*models.Metric, error)
}

func NewTopicDao() TopicDao {
	return scylla.NewTopicScyllaDao()
}
