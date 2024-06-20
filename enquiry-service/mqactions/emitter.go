package mqactions

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Emitter struct {
	connection *amqp.Connection
	exchange   string
	routingKey string
}

func (e *Emitter) setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	return e.declareExchange(channel)
}

func (e *Emitter) declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		e.exchange, // name of exchange
		"direct",   // type
		true,
		false,
		false,
		false,
		nil,
	)
}

func (e *Emitter) Emit(event string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	log.Printf("Pushing to exchange %s with routing key %s", e.exchange, e.routingKey)

	err = channel.Publish(
		e.exchange,
		e.routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func NewEmitter(conn *amqp.Connection, exchange, routingKey string) (*Emitter, error) {
	emitter := &Emitter{
		connection: conn,
		exchange:   exchange,
		routingKey: routingKey,
	}
	err := emitter.setup()

	if err != nil {
		return nil, err
	}

	return emitter, nil
}
