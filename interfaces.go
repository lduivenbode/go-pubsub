package pubsub

type Queue interface {
	Name() string
}

type Subscription interface {
	Queue() Queue
}

type Message interface {
	Headers() map[string]interface{}
	Body() []byte
	Ack() error
	Reject(requeue bool) error
}

type Provider interface {
	Subscribe(s Subscription) (<-chan Message, error)
	Close() error
}

type CloseError interface {
	Close() bool
}

type CancelError interface {
	Cancel() bool
}
