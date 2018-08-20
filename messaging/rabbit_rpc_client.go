package messaging

import (
	"github.com/streadway/amqp"
	"math/rand"
	"github.com/gobuffalo/envy"
	"encoding/json"
)
func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

type rpcAction string

type rpcMessage struct {
	Action rpcAction
	Payload map[string]interface{}
}

type rpcClient struct {
	forever  chan bool
	callback CallBackFunc
	qName    RabbitQueue
}

func NewRpcCient(q_name RabbitQueue, callback CallBackFunc) *rpcClient {
	f := make(chan bool)
	return &rpcClient{f, callback, q_name}
}

type CallBackFunc func(rpcMessage) interface{}

type RabbitQueue string
type RabbitRoutingKey string

const (
	USER_RPC_QUEUE    RabbitQueue = "USER_RPC_QUEUE"
	COMPANY_RPI_QUEUE RabbitQueue = "COMPANY_RPC_QUEUE"
)

func (client rpcClient) Request(rpcm rpcMessage) (res int, err error) {
	var rabbit_mq_uri = envy.Get("RABBITMQ_BIGWIG_URL", "amqp://guest:guest@localhost:5672/")
	conn, err := amqp.Dial(rabbit_mq_uri)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	corrId := randomString(32)

	request, err := json.Marshal(&rpcm)
	if err != nil {
		panic(err)
	}
	err = ch.Publish(
		"",          // exchange
		string(client.qName), // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          request,
		})
	failOnError(err, "Failed to publish a message")

	for d := range msgs {
		if corrId == d.CorrelationId {
			response := &rpcMessage{}

			err := json.Unmarshal(d.Body, &response)
			failOnError(err, "Failed to convert body to integer")

			client.callback(*response)
			break
		}
	}

	return
}


