package models

type TopicProperties struct {
	IndexName string
	StartFrom string
}

type Topic struct {
	ID         string
	Name       string
	Properties TopicProperties
}
