package daos

import (
	"github.com/dwsb/projetomiddleware/distribution/models"
	"github.com/dwsb/projetomiddleware/redis"
)

type SessionDao interface {
	AddMember(*models.Session) error
	ListMembers(string) ([]*models.Session, error)
	RemoveMember(*models.Session) error
}

func NewSessionDao() SessionDao {
	return redis.NewSessionRedisDao()
}
