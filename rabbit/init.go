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
