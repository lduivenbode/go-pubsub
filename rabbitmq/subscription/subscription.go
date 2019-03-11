package subscription

import (
	"github.com/streadway/amqp"
	"fmt"
	"github.com/lduivenbode/go-pubsub/rabbitmq/exchange"
	"github.com/lduivenbode/go-pubsub"
)

type Subscription struct {
	exchange *exchange.Exchange
	queue pubsub.Queue
	topics []string
	autoAck bool
	exclusive bool
	noLocal bool
	noWait bool
	args amqp.Table

}

func New(q pubsub.Queue, options ...func(*Subscription) error) (*Subscription, error) {
	var err error
	s := Subscription{}

	s.exchange = exchange.New()

	s.queue = q

	// Process configuration options to override defaults
	for _, option := range options {
		err = option(&s)
		if err != nil {
			return nil, fmt.Errorf("failed to initialise provider: %s", err)
		}
	}

	return &s, nil
}

func (s *Subscription) Exchange() *exchange.Exchange {
	return s.exchange
}

func (s *Subscription) Queue() pubsub.Queue {
	return s.queue
}

func (s *Subscription) Topics() []string {
	return s.topics
}

func (s *Subscription) AutoAck() bool {
	return s.autoAck
}

func (s *Subscription) Exclusive() bool {
	return s.exclusive
}

func (s *Subscription) NoLocal() bool {
	return s.noLocal
}

func (s *Subscription) NoWait() bool {
	return s.noWait
}

func (s *Subscription) Args() amqp.Table {
	return s.args
}

// Topics sets the topics that this queue should subscribe to
func Topics(topics []string) func(s *Subscription) error {
	return func(s *Subscription) error {
		s.topics = topics
		return nil
	}
}

//func Exclusive(s *Subscription, exclusive bool) error {
//	s.exclusive = exclusive
//	return nil
//}
//
//// TopicExchange sets the name of the topic to use for this exchange (default: 'topic')
//func NoWait(s *Subscription, noWait bool) error {
//	s.noWait = noWait
//	return nil
//}
//
//// Args sets the AMQP options for the declaring of the queue (default nil)
//func Args(s *Subscription, args amqp.Table) error {
//	s.args = args
//	return nil
//}
//
