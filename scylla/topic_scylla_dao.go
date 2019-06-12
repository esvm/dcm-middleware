package scylla

import (
	"fmt"
	"os"
	"strings"

	"github.com/esvm/dcm-middleware/distribution/models"
	"github.com/gocql/gocql"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type TopicScyllaDao struct {
	session *gocql.Session
}

const (
	Begin = "BEGIN"
	End   = "END"

	defaultKeyspace     = "broker"
	defaultHost         = "127.0.0.1"
	createTableTemplate = `CREATE TABLE broker.%s (
		id text PRIMARY KEY,
		average float,
		median float,
		variance float,
		standartDeviation float,
		mode float,
		created_at timestamp, 
	);`
)

func orderBy(indexStart string) string {
	switch indexStart {
	case End:
		return "MAX"
	default:
		return "MIN"
	}
}

func NewTopicScyllaDao() *TopicScyllaDao {
	host := os.Getenv("SCYLLA_HOST")
	if host == "" {
		log.Warn("Environment not setted. Using default host: ", defaultHost)
		host = defaultHost
	}
	hosts := strings.Split(host, ";")

	conf := gocql.NewCluster(hosts...)
	conf.RetryPolicy = &gocql.ExponentialBackoffRetryPolicy{NumRetries: 3}
	conf.Consistency = gocql.LocalQuorum

	s, err := conf.CreateSession()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
		return nil
	}

	return &TopicScyllaDao{session: s}
}

func (t *TopicScyllaDao) CreateTopic(topic *models.Topic) (*models.Topic, error) {
	var IDValue, nameValue string
	qCql := t.session.Query(fmt.Sprintf("SELECT id, name FROM broker.topics WHERE id = '%s'", topic.ID))
	err := qCql.Scan(&IDValue, &nameValue)
	if err == nil {
		return &models.Topic{
			ID:   IDValue,
			Name: nameValue,
		}, nil
	}

	log.Debug("create topic: ", err)

	topic.ID = strings.ToLower(topic.Name)

	qCql = t.session.Query(fmt.Sprintf("INSERT INTO broker.topics (id, name) VALUES ('%s', '%s')", topic.ID, topic.Name))
	if err := qCql.Exec(); err != nil {
		log.Debug(err)
		return nil, err
	}

	qCql.Release()

	qCql = t.session.Query(fmt.Sprintf(createTableTemplate, topic.ID))
	if err := qCql.Exec(); err != nil {
		log.Debug(err)
		return nil, err
	}

	qCql.Release()
	return topic, nil
}

func (t *TopicScyllaDao) InsertMessage(topicID string, metrics *models.Metric) error {
	qCql := t.session.Query(fmt.Sprintf(
		"INSERT INTO broker.%s (id, average, median, variance, standartDeviation, mode, created_at) values ('%s', %f, %f, %f, %f, %f, toUnixTimestamp(now()))",
		topicID,
		uuid.NewV4().String(),
		metrics.Average,
		metrics.Median,
		metrics.Variance,
		metrics.StandardDeviation,
		metrics.Mode))

	if err := qCql.Exec(); err != nil {
		log.Debug(err)
		return err
	}

	qCql.Release()
	return nil
}

func (t *TopicScyllaDao) InsertIndex(topicID, indexName, indexStart string) (*models.Index, error) {
	var tokenID string

	qCql := t.session.Query(fmt.Sprintf("SELECT %s(token(id)) FROM broker.%s", orderBy(indexStart), topicID))
	if err := qCql.Scan(&tokenID); err != nil && err != gocql.ErrNotFound {
		log.Debug("xablau", err)
		return nil, err
	}

	var value string
	if tokenID != "0" {
		qCql = t.session.Query(fmt.Sprintf("SELECT id FROM broker.%s WHERE token(id) = %s", topicID, tokenID))
		if err := qCql.Scan(&value); err != nil {
			log.Debug(err)
			return nil, err
		}
	}

	qCql = t.session.Query(fmt.Sprintf(
		"INSERT INTO broker.indexes (name, value) values ('%s', '%s')",
		indexName,
		value))

	if err := qCql.Exec(); err != nil {
		log.Debug(err)
		return nil, err
	}

	qCql.Release()
	return &models.Index{
		Name:  indexName,
		Value: value,
	}, nil
}

func (t *TopicScyllaDao) IncIndex(topicID string, index *models.Index, limit int) (*models.Index, error) {
	if limit <= 0 {
		return index, nil
	}

	qCql := t.session.Query(fmt.Sprintf("SELECT id FROM broker.%s WHERE token(id) > token('%s') LIMIT %d", topicID, index.Value, limit))
	if err := qCql.Scan(&index.Value); err != nil {
		log.Debug(err)
		return nil, err
	}
	log.Info(index.Value)
	qCql = t.session.Query(fmt.Sprintf("UPDATE broker.indexes SET value = '%s' WHERE name = '%s'", index.Value, index.Name))
	if err := qCql.Exec(); err != nil {
		log.Debug(err)
		return nil, err
	}

	qCql.Release()

	return index, nil
}

func (t *TopicScyllaDao) GetIndex(indexName string) (*models.Index, error) {
	var value string
	qCql := t.session.Query(fmt.Sprintf("SELECT value FROM broker.indexes WHERE name = '%s' LIMIT 1", indexName))
	if err := qCql.Scan(&value); err != nil {
		log.Debug(err)
		return nil, err
	}

	return &models.Index{
		Name:  indexName,
		Value: value,
	}, nil
}

func (t *TopicScyllaDao) GetMessage(topicID string, index *models.Index) (*models.Metric, error) {
	var average, median, variance, standardDeviation, mode float32
	qCql := t.session.Query(fmt.Sprintf("SELECT average, median, variance, standartDeviation, mode FROM broker.%s WHERE id = '%s' LIMIT 1", topicID, index.Value))
	if err := qCql.Scan(&average, &median, &variance, &standardDeviation, &mode); err != nil {
		log.Debug(err)
		return nil, err
	}

	return &models.Metric{
		ID:                index.Value,
		Average:           float64(average),
		Median:            float64(median),
		Variance:          float64(variance),
		StandardDeviation: float64(standardDeviation),
		Mode:              float64(mode),
	}, nil
}

func (t *TopicScyllaDao) GetMessages(topicID string, index *models.Index, limit int) ([]*models.Metric, error) {
	var average, median, variance, standardDeviation, mode float32
	iter := t.session.Query(fmt.Sprintf(
		"SELECT average, median, variance, standartDeviation, mode FROM broker.%s WHERE token(id) >= token('%s') LIMIT %d",
		topicID,
		index.Value,
		limit)).Iter()

	var metrics []*models.Metric
	for iter.Scan(&average, &median, &variance, &standardDeviation, &mode) {
		metric := &models.Metric{
			ID:                index.Value,
			Average:           float64(average),
			Median:            float64(median),
			Variance:          float64(variance),
			StandardDeviation: float64(standardDeviation),
			Mode:              float64(mode),
		}

		metrics = append(metrics, metric)
	}

	if err := iter.Close(); err != nil {
		log.Debug(err)
		return nil, err
	}

	return metrics, nil
}
