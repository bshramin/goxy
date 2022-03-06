package rabbit

import "github.com/streadway/amqp"

type RInfo interface {
	DeclareQueue(ch *amqp.Channel, queueName string) (amqp.Queue, error)
	DeclareExchange(ch *amqp.Channel, exchangeName string) error
}
