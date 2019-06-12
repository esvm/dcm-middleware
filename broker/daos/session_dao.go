package daos

import (
	"github.com/esvm/dcm-middleware/distribution/models"
	"github.com/esvm/dcm-middleware/redis"
)

type SessionDao interface {
	AddMember(*models.Session) error
	ListMembers(string) ([]*models.Session, error)
	RemoveMember(*models.Session) error
}

func NewSessionDao() SessionDao {
	return redis.NewSessionRedisDao()
}
