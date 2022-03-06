package rabbit

import (
	"github.com/streadway/amqp"
)

// Emitter for publishing AMQP events
type Emitter struct {
	Conn         *amqp.Connection
	ExchangeName string
}

func (e *Emitter) setup(r RInfo) error {
	ch, err := e.Conn.Channel()
	if err != nil {
		panic(err)
	}

	return r.DeclareExchange(ch, e.ExchangeName)
}

// NewEventEmitter returns a new event.Emitter object
// ensuring that the object is initialised, without error
func NewEventEmitter(conn *amqp.Connection, r RInfo, exchangeName string) (Emitter, error) {
	emitter := Emitter{
		Conn:         conn,
		ExchangeName: exchangeName,
	}

	err := emitter.setup(r)
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}
