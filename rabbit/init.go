package rabbit

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func GetConnection(user, password, host string, port int, vHost string) *amqp.Connection {
	connString := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		user,
		password,
		host,
		port,
		vHost,
	)
	conn, err := amqp.Dial(connString)
	if err != nil {
		logrus.Errorf("error while connect to rabbit : %s", err.Error())
		panic(err)
	}
	return conn
}

func DeclareQueue(ch *amqp.Channel, queueName string) (amqp.Queue, error) {
	return ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		true,      // exclusive
		false,     // no-wait
		nil,       // arguments
	)
}

func DeclareExchange(ch *amqp.Channel, exchangeName string) error {
	return ch.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
}
