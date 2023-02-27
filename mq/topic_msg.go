package mq

type TopicMessage struct {
	Topic       interface{} `json:"topic"`
	ProjectId   string      `json:"project_id"`
	From        string      `json:"from"`
	Timestamp   int64       `json:"timestamp"`
	Data        interface{} `json:"data"`
	LogId       string      `json:"log_id"`
	MsgId       string      `json:"msg_id"`
	OrderingKey string      `json:"ordering_key"`
}
