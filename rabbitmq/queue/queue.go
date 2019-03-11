package queue

import (
	"github.com/streadway/amqp"
	"fmt"
)

// Queue represents settings for declaring RabbitMQ queues
type Queue struct {
	name string
	durable bool
	autoDelete bool
	exclusive bool
	noWait bool
	args amqp.Table
}

func New(target string, options ...func(*Queue) error) (*Queue, error) {
	q := Queue{}
	q.name = target
	q.durable = true

	// Process configuration options to override defaults
	for _, option := range options {
		err := option(&q)
		if err != nil {
			return nil, fmt.Errorf("failed to initialise provider: %s", err)
		}
	}

	return &q, nil
}

// Name returns the name of the queue
func (q *Queue) Name() string {
	return q.name
}

func (q *Queue) SetName(name string) {
	q.name = name
}

// Durable returns if the queue should survive if all the consumers disconnect
func (q *Queue) Durable() bool {
	return q.durable
}

// AutoDelete returns if the queue should be removed when there are no longer any bindings or consumers
func (q *Queue) AutoDelete() bool {
	return q.autoDelete
}

// Exclusive returns whether the queue should be restricted to the creator
func (q *Queue) Exclusive() bool {
	return q.exclusive
}
// NoWait returns if RabbitMQ should wait for the queue declaration to complete before returning
func (q *Queue) NoWait() bool {
	return q.noWait
}

// Args returns additional arguments used in declaration of the queue
func (q *Queue) Args() amqp.Table {
	return q.args
}

// Durable sets whether the queue is durable or not
func Durable(durable bool) func(q *Queue) error {
	return func (q *Queue) error {
		q.durable = durable
		return nil
	}
}

// AutoDelete sets whether the queue should be deleted when there are no further bindings or consumers
func AutoDelete(autoDelete bool) func(q *Queue) error {
	return func (q *Queue) error {
		q.autoDelete = autoDelete
		return nil
	}
}

// Exclusive returns whether the queue should be restricted to the creator
func Exclusive(q *Queue, exclusive bool) func(q *Queue) error {
	return func (q *Queue) error {
		q.exclusive = exclusive
		return nil
	}
}

// NoWait returns if RabbitMQ should wait for the queue declaration to complete before returning
func NoWait(q *Queue, noWait bool) func(q *Queue) error {
	return func (q *Queue) error {
		q.noWait = noWait
		return nil
	}
}

// Args returns additional arguments used in declaration of the queue
func Args(q *Queue, args amqp.Table) func(q *Queue) error {
	return func (q *Queue) error {
		q.args = args
		return nil
	}
}
