package models

const (
	CreateTopicOperation = "create_topic"
	PublishOperation     = "publish"
	ConsumeOperation     = "consume"
	ResetIndexOperation  = "reset_index"
)

type Request interface {
	Operation() string
}

type CreateTopicRequest struct {
	Name       string
	Properties TopicProperties
}

func NewCreateTopicRequest(name string, properties TopicProperties) *CreateTopicRequest {
	return &CreateTopicRequest{
		Name:       name,
		Properties: properties,
	}
}

func (req *CreateTopicRequest) Operation() string {
	return CreateTopicOperation
}

type PublishRequest struct {
	TopicID string
	Metrics *Metric
}

func NewPublishRequest(topicID string, metrics *Metric) *PublishRequest {
	return &PublishRequest{
		TopicID: topicID,
		Metrics: metrics,
	}
}

func (req *PublishRequest) Operation() string {
	return PublishOperation
}

type ConsumeRequest struct {
	TopicID   string
	IndexName string
}

func NewConsumeRequest(topicID, indexName string) *ConsumeRequest {
	return &ConsumeRequest{
		TopicID:   topicID,
		IndexName: indexName,
	}
}

func (req *ConsumeRequest) Operation() string {
	return ConsumeOperation
}

type ResetIndexRequest struct {
	TopicName string
	IndexName string
}

func NewResetIndexRequest(topicName, indexName string) *ResetIndexRequest {
	return &ResetIndexRequest{
		TopicName: topicName,
		IndexName: indexName,
	}
}

func (req *ResetIndexRequest) Operation() string {
	return ResetIndexOperation
}
