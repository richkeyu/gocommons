package mq

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	application "github.com/richkeyu/gocommons/dispatch"
	"github.com/richkeyu/gocommons/plog"

	"cloud.google.com/go/pubsub"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/api/option"
)

type Consumer struct {
	Name         string
	ProjectId    string
	ConsumerId   string
	Path         string
	ErrSendEmail map[string]string
	Usage        string
	TimeOut      time.Duration //默认10秒
	ThreadNum    int           //是否异步并行处理 默认 4
	Action       func(msg *TopicMessage) error
	stopChan     chan error
	isStopped    bool
	client       *pubsub.Client
	waitGroup    sync.WaitGroup // 未使用
}

type ConsumerOption struct {
	Name         string `yaml:"name"`
	ConfPath     string `yaml:"path"` // host:port
	ConsumerId   string `yaml:"consumer_id"`
	ProjectId    string `yaml:"project_id"`
	Usage        string `yaml:"usage"`
	ThreadNum    int    `yaml:"thread_num"` //是否异步并行处理 默认 4
	ErrSendEmail map[string]string
	Action       func(msg *TopicMessage) error
}

func NewConsumer(consumer *ConsumerOption) *Consumer {
	c := &Consumer{
		Name:         consumer.Name,
		ProjectId:    consumer.ProjectId,
		ConsumerId:   consumer.ConsumerId,
		Path:         consumer.ConfPath,
		ErrSendEmail: consumer.ErrSendEmail,
		Usage:        consumer.Usage,
		ThreadNum:    consumer.ThreadNum,
		Action: func(msg *TopicMessage) error {
			return consumer.Action(msg)
		},
	}
	return c
}

func (c *Consumer) initConfig() {
	if c.ThreadNum < 1 || c.ThreadNum > 10 {
		c.ThreadNum = 4
	}
}

func (c *Consumer) start() (err error) {
	c.initConfig()
	defer func() {
		c.isStopped = true
		close(c.stopChan)
	}()
	c.stopChan = make(chan error)

	ctx, cancel := context.WithCancel(context.Background())
	client, err := pubsub.NewClient(ctx, c.ProjectId, option.WithCredentialsFile(c.Path))
	if err != nil {
		plog.Errorf(nil, "consumer NewClient, error:%v", err)
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	application.OnExit(cancel)

	c.client = client
	sub := client.Subscription(c.ConsumerId)
	sub.ReceiveSettings.Synchronous = false
	sub.ReceiveSettings.NumGoroutines = c.ThreadNum
	sub.ReceiveSettings.MaxOutstandingMessages = 100 // todo 暂定100,看情况是否支持配置
	plog.Infof(nil, "projectId:%s||consumerId:%s||start run", c.ProjectId, c.ConsumerId)
	err = sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		handleErr := c.doAction(c.ProjectId, c.ConsumerId, newTopicMessage4Consumer(msg), func() {
		})
		if handleErr != nil {
			msg.Nack()
			return
		}
		msg.Ack()
	})
	if err != nil {
		plog.Errorf(nil, "projectId:%s||consumerId:%s||receive error: %v", c.ProjectId, c.ConsumerId, err)
		return err
	}

	plog.Infof(nil, "projectId:%s||consumerId:%s||stop run", c.ProjectId, c.ConsumerId)
	time.Sleep(time.Second * 1) // 为了将上行结束日志记录到磁盘，所以让消费协程等待1s退出
	return
}

func (c *Consumer) close() {
	err := c.client.Close()
	if err != nil {
		plog.Error(nil, fmt.Sprintf("close consumer, projectId:%s consumerId:%s err error:%v", c.ProjectId, c.ConsumerId, err.Error()))
	}
}

func (c *Consumer) doAction(projectId, consumerId string, msgObj *TopicMessage, done func()) error {
	b, _ := jsoniter.Marshal(msgObj)
	var err error
	defer done()
	defer func() {
		if rErr := recover(); rErr != nil {
			logStr := fmt.Sprintf("c defer panic err projectId:%s consumerId:%s value:%s "+
				" error:%v %v", projectId, consumerId, string(b), rErr, string(debug.Stack()))
			plog.Error(nil, logStr)
		}
	}()
	err = c.Action(msgObj)
	if err != nil {
		plog.Error(nil, fmt.Sprintf("c projectId:%s consumerId:%s Action err value:%s error:%v", projectId, consumerId, string(b), err.Error()))
	}
	return err
}

func newTopicMessage4Consumer(msg *pubsub.Message) *TopicMessage {
	data := &TopicMessage{}
	_ = jsoniter.Unmarshal(msg.Data, data)
	data.MsgId = msg.ID

	return data
}
