package broker

import (
	"github.com/dwsb/projetomiddleware/broker/daos"
	"github.com/dwsb/projetomiddleware/distribution/models"
)

type SessionMgmt struct {
	dao daos.SessionDao
}

func NewSessionMgmt() *SessionMgmt {
	return &SessionMgmt{dao: daos.NewSessionDao()}
}

func (s *SessionMgmt) Subscribe(session *models.Session) error {
	return s.dao.AddMember(session)
}

func (s *SessionMgmt) GetMembers(topicID string) ([]*models.Session, error) {
	return s.dao.ListMembers(topicID)
}
