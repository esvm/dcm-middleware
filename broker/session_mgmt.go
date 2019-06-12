package broker

import (
	"github.com/esvm/dcm-middleware/broker/daos"
	"github.com/esvm/dcm-middleware/distribution/models"
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
