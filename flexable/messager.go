package flexable
import "github.com/streadway/amqp"

type hub struct {
	channels map[string]chan interface{}

	sender chan interface{}
}

func (h hub) register (name string, channel chan interface{}) {
	h.channels[name] =  channel
}

func (h hub) listen () {
	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
}