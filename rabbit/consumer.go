package rabbit

import (
	"github.com/streadway/amqp"
)

// Consumer for receiving AMPQ events
type Consumer struct {
	Conn         *amqp.Connection
	QueueName    string
	ExchangeName string
}

func (consumer *Consumer) setup() error {
	ch, err := consumer.Conn.Channel()
	if err != nil {
		return err
	}

	err = DeclareExchange(ch, consumer.ExchangeName)
	if err != nil {
		return err
	}

	q, err := DeclareQueue(ch, consumer.QueueName)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		q.Name,                // queue name
		"",                    // routing key
		consumer.ExchangeName, // exchange
		false,
		nil,
	)
	return err
}

// NewConsumer returns a new Consumer
func NewConsumer(conn *amqp.Connection, exchangeName string, queueName string) (Consumer, error) {
	consumer := Consumer{
		Conn:         conn,
		QueueName:    queueName,
		ExchangeName: exchangeName,
	}
	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}
