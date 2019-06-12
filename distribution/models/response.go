package models

type Response interface {
	Error() string
	KeepAlive() bool
}

type CreateTopicResponse struct {
	TopicResult *Topic
	Err         string
	Alive       bool
}

func (res *CreateTopicResponse) Error() string {
	return res.Err
}

func (res *CreateTopicResponse) KeepAlive() bool {
	return res.Alive
}

type PublishResponse struct {
	Err   string
	Alive bool
}

func (res *PublishResponse) Error() string {
	return res.Err
}

func (res *PublishResponse) KeepAlive() bool {
	return res.Alive
}

type ConsumeResponse struct {
	MessageResult *Metric
	Err           string
	Alive         bool
}

func (res *ConsumeResponse) Error() string {
	return res.Err
}

func (res *ConsumeResponse) KeepAlive() bool {
	return res.Alive
}

type ResetIndexResponse struct {
	Err   string
	Alive bool
}

func (res *ResetIndexResponse) Error() string {
	return res.Err
}

func (res *ResetIndexResponse) KeepAlive() bool {
	return res.Alive
}
