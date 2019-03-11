package rabbitmq

import "github.com/streadway/amqp"

type RMQClose struct {
	AMQPError *amqp.Error
}

func (ce RMQClose) Close() bool {
	return true
}

func (ce RMQClose) Error() string {
	if ce.AMQPError != nil {
		return ce.AMQPError.Error()
	}
	return ""
}

func FromAMQPError(ae *amqp.Error) RMQClose {
	e := RMQClose{ae}
	return e
}

type RMQCancel struct {
	ConsumerTag string
}

func (ce RMQCancel) Cancel() bool {
	return true
}

func (ce RMQCancel) Error() string {
	return ce.ConsumerTag
}
