package mq

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

func TestPublish(t *testing.T) {

	_ = InitProducerConf(&ProducerConf{
		Path:      "/Users/yangfan/work/im30/password/kop-pay-003136049374.json",
		ProjectID: "kop-pay",
		Topic:     "test_mq_queue",
	})
	for i := 260; i < 290; i++ {
		id, err := Publish(context.WithValue(context.Background(), "exec", "test"), "kop-pay", "test_mq_queue", fmt.Sprintf("测试消费者ddddddd%d", i))
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(id)
	}
}

func TestConsumer(t *testing.T) {
	Go(RegisterCommand)
}

func PullMsgs(projectID, subID string) error {
	//	// projectID := "my-project-id"
	//	// subID := "my-sub"
	ctx, cancel := context.WithCancel(context.Background())
	client, err := pubsub.NewClient(ctx, projectID, option.WithCredentialsFile("/Users/yangfan/work/im30/password/kop-pay-003136049374.json"))
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	sub := client.Subscription(subID)

	sub.ReceiveSettings.MaxOutstandingMessages = 1
	sub.ReceiveSettings.NumGoroutines = 1
	// Receive messages for 10 seconds, which simplifies testing.
	// Comment this out in production, since `Receive` should
	// be used as a long running operation.
	done := make(chan bool)
	var received int32
	go func() {
		fmt.Println("success")
		err = sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
			time.Sleep(time.Second * 3)
			if msg != nil {
				fmt.Println(fmt.Sprintf("Got time:%s, message: %#v ", time.Now().Format("2006-01-02 15:04:05"), string(msg.Data)))
			} else {
				fmt.Println(fmt.Sprintf("Got time:%s, message: %#v ", time.Now().Format("2006-01-02 15:04:05"), string(msg.Data)))
			}

			//e := client.Close()
			//if e != nil {
			//	fmt.Println(e.Error())
			//}

			//r := msg.AckWithResult()
			//status, err := r.Get(ctx)
			//if err != nil {
			//	fmt.Println(fmt.Sprintf("MessageID: %s failed when calling result.Get: %v", msg.ID, err))
			//}
			//fmt.Println(fmt.Sprintf("Message successfully status: %d", status))
			msg.Ack()

		})
		if err != nil {
			fmt.Println(fmt.Sprintf("sub.Receive: %v", err))
		} else {
			fmt.Println(fmt.Sprintf("sub. stop"))
		}

		time.Sleep(time.Second * 1)
		done <- true
		return
	}()

	for {
		fmt.Println(fmt.Sprintf("Received %d messages\n", received))
		time.Sleep(time.Second * 4)
		fmt.Println(fmt.Sprintf("Received %d cancel\n", received))
		cancel()
		break
	}

	<-done
	return nil
}

func RegisterCommand(command Run) {
	command.AddConsumer(&ConsumerOption{
		Name:         `test_common_queue`,
		ConfPath:     `/Users/yangfan/work/im30/password/kop-pay-003136049374.json`,
		Usage:        "测试用法",
		ErrSendEmail: map[string]string{`dongbufan@soyoung.com`: "董不凡"},
		ProjectId:    "kop-pay",
		ConsumerId:   "test_mq_common",
		ThreadNum:    1,
		Action:       SendMsgTaskConsumer{}.Do,
	})
}

type SendMsgTaskConsumer struct {
}

func (t SendMsgTaskConsumer) Do(msg *TopicMessage) (err error) {
	time.Sleep(1 * time.Second)
	fmt.Println("我是测试消费者", msg.Data.(string))

	return nil
}
