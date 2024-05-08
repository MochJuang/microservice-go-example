package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Emitter struct {
	connection *amqp.Connection
	Channel    *amqp.Channel
}

func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{connection: conn}
	err := emitter.setup()
	return emitter, err
}

func (e *Emitter) setup() error {
	var err error
	e.Channel, err = e.connection.Channel()
	if err != nil {
		return err
	}
	return nil
}

// severity as routing key
// event as payload
func (e Emitter) Push(key string, data string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	log.Println("push to rabbitmq event:", key)

	err = channel.Publish(
		EXCHANGE_LOGS_TOPIC,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(data),
		},
	)

	if err != nil {
		return err
	}

	return nil
}
