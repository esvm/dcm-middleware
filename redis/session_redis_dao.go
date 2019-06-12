package redis

import (
	"os"

	"github.com/dwsb/projetomiddleware/distribution/models"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

type SessionRedisDao struct {
	client *redis.Client
}

const defaultHost = "localhost"
const defaultPort = "6379"

func NewSessionRedisDao() *SessionRedisDao {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")

	if host == "" {
		log.Warn("Environment not setted. Using default host: ", defaultHost)
		host = defaultHost
	}

	if port == "" {
		log.Warn("Environment not setted. Using default port: ", defaultPort)
		port = defaultPort
	}

	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	return &SessionRedisDao{client: client}
}

func (s *SessionRedisDao) AddMember(session *models.Session) error {
	return s.client.SAdd(session.TopicID, session).Err()
}

func (s *SessionRedisDao) ListMembers(topicID string) ([]*models.Session, error) {
	var members []*models.Session
	err := s.client.SMembers(topicID).ScanSlice(&members)

	return members, err
}

func (s *SessionRedisDao) RemoveMember(session *models.Session) error {
	return s.client.SRem(session.TopicID, session).Err()
}
