package main

import (
"context"
"fmt"
"log"

"github.com/streadway/amqp"
)

//Rabbitmq 初始化rabbitmq连接
type Rabbitmq struct {
	conn *amqp.Connection
	err  error
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}


//New 开始创建一个新的rabitmq连接
func New(ip string) (*Rabbitmq, error) {
	amqps := fmt.Sprintf("amqp://guest:guest@%s:5672/", ip)
	conn, err := amqp.Dial(amqps)
	if err != nil {
		return nil, err
	}
	rabbitmq := &Rabbitmq{
		conn: conn,
	}
	return rabbitmq, nil
}

//CreateQueue 创建一个queue队列
func (rabbitmq *Rabbitmq) CreateQueue(id string) error {
	ch, err := rabbitmq.conn.Channel()
	defer ch.Close()
	if err != nil {
		return err
	}
	_, err = ch.QueueDeclare(
		id,    // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}
	return nil
}

//DeleteQueue 删除一个queue队列
func (rabbitmq *Rabbitmq) DeleteQueue(id string) error {
	ch, err := rabbitmq.conn.Channel()
	defer ch.Close()
	if err != nil {
		return err
	}
	_, err = ch.QueueDelete(
		id,    // name
		false, // IfUnused
		false, // ifEmpty
		true,  // noWait
	)
	if err != nil {
		return err
	}
	return nil
}

//PublishQueue 上传消息到queue队列中
func (rabbitmq *Rabbitmq) PublishQueue(id string, body string) error {
	ch, err := rabbitmq.conn.Channel()
	defer ch.Close()
	if err != nil {
		return err
	}
	err = ch.Publish(
		"",    // exchange
		id,    // routing key
		false, // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	if err != nil {
		return err
	}
	return nil
}

//ConsumeQueue 从队列中取出数据并且消费
func (rabbitmq *Rabbitmq) ConsumeQueue(ctx context.Context, id string) error {
	ch, err := rabbitmq.conn.Channel()
	if err != nil {
		return err
	}
	err = ch.Qos(
		3,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		id,    // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			//消费数据
			fmt.Println(string(d.Body))
			//标记消费
			d.Ack(false)
		}
	}()
	select {
	case <-ctx.Done():
		fmt.Println("任务结束")
		return nil
	}
	return nil
}


//GetReadyCount 统计正在队列中准备且还未消费的数据
func (rabbitmq *Rabbitmq) GetReadyCount(id string) (int, error) {
	count := 0
	ch, err := rabbitmq.conn.Channel()
	defer ch.Close()
	if err != nil {
		return count, err
	}
	state, err := ch.QueueInspect(id)
	if err != nil {
		return count, err
	}
	return state.Messages, nil
}

//GetConsumCount 获取到队列中正在消费的数据，这里指的是正在有多少数据被消费
func (rabbitmq *Rabbitmq) GetConsumCount(id string) (int, error) {
	count := 0
	ch, err := rabbitmq.conn.Channel()
	defer ch.Close()
	if err != nil {
		return count, err
	}
	state, err := ch.QueueInspect(id)
	if err != nil {
		return count, err
	}
	return state.Consumers, nil
}

//ClearQueue 清理队列
func (rabbitmq *Rabbitmq) ClearQueue(id string) (string, error) {
	ch, err := rabbitmq.conn.Channel()
	defer ch.Close()
	if err != nil {
		return "", err
	}
	_, err = ch.QueuePurge(id, false)
	if err != nil {
		return "", err
	}
	return "Delete queue success", nil
}
