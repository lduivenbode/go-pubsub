package exchange

import "github.com/streadway/amqp"

// Exchange represents settings for declaring RabbitMQ exchanges
type Exchange struct {
	name       string
	durable    bool
	autoDelete bool
	noWait     bool
	args       amqp.Table
}

func New() *Exchange {
	e := Exchange{}
	e.name = "topic"
	e.durable = true
	return &e
}

// Name returns the name of the exchange
func (e *Exchange) Name() string {
	return e.name
}

// Durable returns if the exchange should survive if all the consumers disconnect
func (e *Exchange) Durable() bool {
	return e.durable
}

// AutoDelete returns if the exchange should be removed when there are no longer any bindings or consumers
func (e *Exchange) AutoDelete() bool {
	return e.autoDelete
}

// NoWait returns if RabbitMQ should wait for the exchange declaration to complete before returning
func (e *Exchange) NoWait() bool {
	return e.noWait
}

// Args returns additional arguments used in declaration of the exchange
func (e *Exchange) Args() amqp.Table {
	return e.args
}

func Name(name string) func(e *Exchange) error {
	return func(e *Exchange) error {
		e.name = name
		return nil
	}
}

// Durable sets whether the queue is durable or not
func Durable(durable bool) func(e *Exchange) error {
	return func(e *Exchange) error {
		e.durable = durable
		return nil
	}
}

// AutoDelete sets whether the queue should be deleted when there are no further bindings or consumers
func AutoDelete(autoDelete bool) func(e *Exchange) error {
	return func(e *Exchange) error {
		e.autoDelete = autoDelete
		return nil
	}
}

// NoWait returns if RabbitMQ should wait for the queue declaration to complete before returning
func NoWait(noWait bool) func(e *Exchange) error {
	return func(e *Exchange) error {
		e.noWait = noWait
		return nil
	}
}

// Args returns additional arguments used in declaration of the queue
func Args(args amqp.Table) func(e *Exchange) error {
	return func(e *Exchange) error {
		e.args = args
		return nil
	}
}
