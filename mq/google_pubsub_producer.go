package mq

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/richkeyu/gocommons/trace"

	"cloud.google.com/go/pubsub"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/api/option"
)

var producers = map[string]*pubsub.Topic{}
var mu sync.Mutex
var confMu sync.Mutex

var producerConf = map[string]*ProducerConf{}

type ProducerConf struct {
	Path      string `yaml:"path"`
	ProjectID string `yaml:"project_id"`
	Topic     string `yaml:"topic"`
}

func getProducerConf(key string) *ProducerConf {
	conf, ok := producerConf[key]
	if ok {
		return conf
	}

	return nil
}

// 初始化生产者conf
func InitProducerConf(conf *ProducerConf) error {
	confMu.Lock()
	defer confMu.Unlock()
	if conf == nil {
		return fmt.Errorf("initProducerConf conf non-existent")
	}

	key := fmt.Sprintf("%s_%s", conf.ProjectID, conf.Topic)
	producerConf[key] = conf

	return nil
}

func createProducer(projectId string, topic string) (*pubsub.Topic, error) {
	key := fmt.Sprintf("%s_%s", projectId, topic)
	if producers == nil {
		producers = map[string]*pubsub.Topic{}
	}
	producer, ok := producers[key]
	if ok {
		return producer, nil
	}
	mu.Lock()
	defer mu.Unlock()
	conf := getProducerConf(key)
	if conf == nil {
		return nil, fmt.Errorf("createProducer conf non-existent")
	}
	if conf.Path == "" {
		return nil, fmt.Errorf("createProducer conf.Path is empty")
	}

	if conf.ProjectID == "" || conf.ProjectID != projectId {
		return nil, fmt.Errorf("createProducer conf.ProjectID is empty")
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, conf.ProjectID, option.WithCredentialsFile(conf.Path))
	if err != nil {
		return nil, fmt.Errorf("pubsub: NewClient: %v", err)
	}

	t := client.Topic(topic)
	producers[key] = t
	return t, nil
}

func Publish(ctx context.Context, projectId, topicId string, msg interface{}) (msgId string, err error) {
	key := fmt.Sprintf("%s_%s", projectId, topicId)
	producer, _ := producers[key]
	if producer == nil {
		producer, err = createProducer(projectId, topicId)
		if err != nil {
			return
		}
	}

	// 获取发送消息调用位置
	_, codePath, codeLine, ok := runtime.Caller(1)
	fileName := ``
	if ok {
		fileName = codePath + `:` + strconv.Itoa(codeLine)
	}
	data := &TopicMessage{
		Topic:     topicId,
		ProjectId: projectId,
		From:      fileName,
		Timestamp: time.Now().Unix(),
		Data:      msg,
		LogId:     trace.GetTraceIdFromContext(ctx),
	}
	b, e := jsoniter.Marshal(data)
	if e != nil {
		return "", e
	}
	newCtx := context.Background()
	message := &pubsub.Message{
		Data: b,
	}
	result := producer.Publish(newCtx, message)
	msgId, err = result.Get(newCtx)
	if err != nil {
		return
	}

	return
}

func PublishOrdering(ctx context.Context, projectId, topicId string, orderingKey string, msg interface{}) (msgId string, err error) {
	key := fmt.Sprintf("%s_%s", projectId, topicId)
	producer, _ := producers[key]
	if producer == nil {
		producer, err = createProducer(projectId, topicId)
		if err != nil {
			return
		}
	}

	// 获取发送消息调用位置
	_, codePath, codeLine, ok := runtime.Caller(1)
	fileName := ``
	if ok {
		fileName = codePath + `:` + strconv.Itoa(codeLine)
	}
	data := &TopicMessage{
		Topic:       topicId,
		ProjectId:   projectId,
		From:        fileName,
		Timestamp:   time.Now().Unix(),
		Data:        msg,
		LogId:       trace.GetTraceIdFromContext(ctx),
		OrderingKey: orderingKey,
	}
	b, e := jsoniter.Marshal(data)
	if e != nil {
		return "", e
	}
	newCtx := context.Background()
	message := &pubsub.Message{
		Data: b,
	}
	if orderingKey != "" {
		producer.EnableMessageOrdering = true
		message.OrderingKey = orderingKey
	}
	result := producer.Publish(newCtx, message)
	msgId, err = result.Get(newCtx)
	if err != nil {
		return
	}

	return
}

//func PullMsgs(projectID, subID string) error {
//	// projectID := "my-project-id"
//	// subID := "my-sub"
//	ctx := context.Background()
//	client, err := pubsub.NewClient(ctx, projectID, option.WithCredentialsFile("/Users/yangfan/work/im30/password/kop-pay-003136049374.json"))
//	if err != nil {
//		return fmt.Errorf("pubsub.NewClient: %v", err)
//	}
//	defer client.Close()
//
//	sub := client.Subscription(subID)
//
//	// Receive messages for 10 seconds, which simplifies testing.
//	// Comment this out in production, since `Receive` should
//	// be used as a long running operation.
//	var received int32
//	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
//		fmt.Println(fmt.Sprintf("Got time:%s, message: %q\n", time.Now().Format("2006-01-02 15:04:05"), string(msg.Data)))
//		//r := msg.AckWithResult()
//		//status, err := r.Get(ctx)
//		//if err != nil {
//		//	fmt.Println(fmt.Sprintf("MessageID: %s failed when calling result.Get: %v", msg.ID, err))
//		//}
//		//fmt.Println(fmt.Sprintf("Message successfully status: %d", status))
//
//		atomic.AddInt32(&received, 1)
//		msg.Ack()
//	})
//	if err != nil {
//		return fmt.Errorf("sub.Receive: %v", err)
//	}
//	fmt.Println(fmt.Sprintf("Received %d messages\n", received))
//	time.Sleep(time.Second * 1)
//
//	return nil
//}
