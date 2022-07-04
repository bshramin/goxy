package rabbit

import (
	"fmt"

	"github.com/streadway/amqp"
)

func GetConnection(user, password, host string, port int, vHost string) (*amqp.Connection, error) {
	connString := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		user,
		password,
		host,
		port,
		vHost,
	)
	conn, err := amqp.Dial(connString)

	return conn, err
}
