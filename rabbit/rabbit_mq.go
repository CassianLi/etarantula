package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

type RabbitMQ struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	exchange     string
	exchangeType string
}

// NewRabbitMQ 创建一个RabbitMQ实例，如果exchange为空则创建一个默认的exchange。heartbeat为心跳间隔，尽量不要设置的太长
func NewRabbitMQ(amqpURI, exchange, exchangeType string, heartbeat time.Duration) (*RabbitMQ, error) {
	mq := &RabbitMQ{
		conn:         nil,
		channel:      nil,
		exchange:     exchange,
		exchangeType: exchangeType,
	}

	var err error
	if heartbeat == 0 {
		heartbeat = time.Second * 10
	}

	mq.conn, err = amqp.DialConfig(amqpURI, amqp.Config{
		Heartbeat: time.Second * 30, // Set the heartbeat interval to 30 seconds
	})
	if err != nil {
		return nil, err
	}

	fmt.Println("Heartbeat interval set to:", heartbeat)

	mq.channel, err = mq.conn.Channel()
	if err != nil {
		return nil, err
	}

	if exchange != "" {
		err = mq.channel.ExchangeDeclare(
			mq.exchange,     // name of the exchange
			mq.exchangeType, // type
			true,            // durable
			false,           // delete when complete
			false,           // internal
			false,           // noWait
			nil,             // arguments
		)
		if err != nil {
			return nil, err
		}
	}

	// 启动心跳发送并记录日志
	go mq.startHeartbeatSender(heartbeat)

	return mq, nil
}

// Close closes the RabbitMQ connection and channel.
func (mq *RabbitMQ) Close() {
	err := mq.channel.Close()
	if err != nil {
		fmt.Println("failed to close channel:", err)
	}
	err = mq.conn.Close()
	if err != nil {
		fmt.Println("failed to close connection:", err)
	}
}

func (mq *RabbitMQ) startHeartbeatSender(heartbeatInterval time.Duration) {
	go func() {
		heartbeatTicker := time.NewTicker(heartbeatInterval)
		defer heartbeatTicker.Stop()

		for range heartbeatTicker.C {
			log.Println("Sending heartbeat...")
		}
	}()
}

// Publish 发布消息到指定的queue，如果没有设置exchange则使用默认的exchange。如果队列不存在则会创建。
func (mq *RabbitMQ) Publish(queue string, body []byte) error {
	if err := mq.channel.Publish(
		mq.exchange, // exchange
		queue,       // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		}); err != nil {
		return err
	}
	return nil
}

// Consume 消费指定的queue，如果队列不存在则会创建。并通过回调函数处理消息，可在调用时指定是否自动应答，如果不自动应答，则在回调函数处理完消息后需要手动应答。Consume方法会阻塞直到出现错误或者channel关闭。
func (mq *RabbitMQ) Consume(queue string, autoAck bool, callback func(msg string)) error {
	deliveries, err := mq.channel.Consume(
		queue,   // name
		"",      // consumerTag,
		autoAck, // autoAck
		false,   // exclusive
		false,   // noLocal
		false,   // noWait
		nil,     // arguments
	)
	if err != nil {
		return err
	}

	for d := range deliveries {
		msg := string(d.Body)

		// 处理消息
		callback(msg)

		// 手动应答
		if !autoAck {
			if err := d.Ack(false); err != nil {
				fmt.Println("failed to ack:", err)
			}
		}
	}

	return nil
}
