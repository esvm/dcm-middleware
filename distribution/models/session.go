package models

import (
	"encoding/json"
)

type Session struct {
	TopicID   string
	Host      string
	IndexName string
}

func (s *Session) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Session) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}
