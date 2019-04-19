package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"
	"github.com/lduivenbode/go-pubsub"
	"github.com/lduivenbode/go-pubsub/rabbitmq/connection"
	"github.com/lduivenbode/go-pubsub/rabbitmq/queue"
	"github.com/lduivenbode/go-pubsub/rabbitmq/subscription"
)

// Provider represents a RabbitMQ concrete implementation of the pubsub.Provider interface
type Provider struct {
	conn *connection.Connection
	ch   *amqp.Channel

	closeCh  chan pubsub.CloseError
	cancelCh chan pubsub.CancelError

	consumerTag   string
	topicExchange string
	prefetchCount int
	prefetchSize  int
	global        bool
}

// ConsumerTag sets the tag for this consumer (instead of the automatically generated tag)
func ConsumerTag(tag string) func(p *Provider) error {
	return func(p *Provider) error {
		p.consumerTag = tag
		return nil
	}
}

// New - returns a new instance of the RabbitMQ provider
func New(options ...func(*Provider) error) (*Provider, error) {
	var err error

	p := Provider{}

	// Default config
	p.topicExchange = "topic"
	p.prefetchCount = 1
	p.prefetchSize = 0
	p.global = false

	// Process configuration options to override defaults
	for _, option := range options {
		// TODO - don't process NotifyClose until later
		err = option(&p)
		if err != nil {
			return nil, fmt.Errorf("failed to initialise provider: %s", err)
		}
	}

	// Create a new connection if one has been created yet
	if p.conn == nil {
		p.conn, err = connection.New()
		if err != nil {
			return nil, fmt.Errorf("failed to initialise provider: %s", err)
		}
	}

	p.ch, err = p.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to initialise provider: %s", err)
	}

	if err = p.ch.Qos(p.prefetchCount, p.prefetchSize, p.global); err != nil {
		return nil, fmt.Errorf("unable to initialise provider - %s", err)
	}

	return &p, nil
}

// Subscribe - subscribe
func (p *Provider) Subscribe(s *subscription.Subscription) (<-chan *Message, error) {
	q := s.Queue()
	e := s.Exchange()

	if e.Name() != "" && len(s.Topics()) > 0 {
		err := p.ch.ExchangeDeclare(
			e.Name(),       // exchange
			"topic",        // kind
			e.Durable(),    // durable
			e.AutoDelete(), // autoDelete
			false,          // internal - no valid use cases for consumers
			e.NoWait(),     // noWait
			e.Args(),       // args
		)
		if err != nil {
			return nil, fmt.Errorf("unable to subscribe - %s", err)
		}
	}

	// Declare the specified queue (can be empty for to use an
	// auto-generated name - which is why we use q.Name from now on)
	rq, err := p.ch.QueueDeclare(
		q.Name(),
		q.Durable(),
		q.AutoDelete(),
		q.Exclusive(),
		q.NoWait(),
		q.Args(),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to subscribe - %s", err)
	}

	q.SetName(rq.Name)

	for _, t := range s.Topics() {
		if err = p.ch.QueueBind(q.Name(), t, e.Name(), false, nil); err != nil {
			return nil, fmt.Errorf("unable to subscribe - %s", err)
		}
	}

	msgs, err := p.ch.Consume(
		q.Name(),        // queue
		p.consumerTag,   // consumer
		s.AutoAck(),   // autoAck
		s.Exclusive(), // exclusive
		s.NoLocal(),   // noLocal
		s.NoWait(),    // noWait
		s.Args(),      // args
	)
	if err != nil {
		return nil, err
	}

	// Consumer amqp.Delivery messages and wrap them into pubsub.Message compatible objects
	c := make(chan pubsub.Message)
	go func(c chan pubsub.Message, rmqMsgs <-chan amqp.Delivery) {
		for msg := range rmqMsgs {
			m := Message{msg}
			c <- &m
		}
	}(c, msgs)

	return c, nil

}

// Close close the connection to the RabbitMQ server
func (p *Provider) Close() error {
	if err := p.conn.Close(); err != nil {
		return fmt.Errorf("unable to close AMQP connection: %s", err)
	}

	return nil
}

func Connection(c *connection.Connection) func(p *Provider) error {
	return func(p *Provider) error {
		p.conn = c
		return nil
	}
}

func NotifyClose(c chan pubsub.CloseError) func(p *Provider) error {
	return func(p *Provider) error {
		ac := make(chan *amqp.Error)

		go func(c chan pubsub.CloseError, ac chan *amqp.Error) {
			// Pipe converted errors to our listener
			for e := range ac {
				c <- FromAMQPError(e)
			}

			// When the AMQP error channel closes we need to close our channel
			close(c)
		}(c, ac)

		return nil
	}
}

func NotifyCancel(c chan pubsub.CancelError) func(p *Provider) error {
	return func(p *Provider) error {
		ac := make(chan string)

		go func(c chan pubsub.CancelError, ac chan string) {
			// Pipe converted errors to our listener
			for e := range ac {
				c <- RMQCancel{e}
			}

			// When the AMQP error channel closes we need to close our channel
			close(c)
		}(c, ac)

		return nil
	}
}
