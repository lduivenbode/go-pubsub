package rabbitmq

import (
	"github.com/streadway/amqp"
)

// Message - RabbitMQ concrete implementation of the pubsub.Message interface
type Message struct {
	msg amqp.Delivery
}

// Headers - return message headers
func (m *Message) Headers() map[string]interface{} {
	return m.msg.Headers
}

// Body - return message body
func (m *Message) Body() []byte {
	return m.msg.Body
}

// Ack - acknowledge the message
func (m *Message) Ack() error {
	return m.msg.Ack(false)
}

// Reject - reject the message
func (m *Message) Reject(requeue bool) error {
	return m.msg.Reject(requeue)
}
