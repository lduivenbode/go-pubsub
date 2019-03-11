package connection

import (
	"fmt"

	"github.com/streadway/amqp"
)

type Connection struct {
	conn     *amqp.Connection
	host     string
	user     string
	password string
	vHost    string
	port     int16
	ssl      bool
}

// Host sets the host for this consumer's connection
func Host(c *Connection, host string) error {
	c.host = host
	return nil
}

// Port sets the port for this consumer's connection
func Port(c *Connection, port int16) error {
	c.port = port
	return nil
}

// User sets the username for this consumer's connection
func User(c *Connection, user string) error {
	c.user = user
	return nil
}

// Password sets the password for this consumer's connection
func Password(c *Connection, password string) error {
	c.password = password
	return nil
}

// VHost sets the vHost for this consumer's connection
func VHost(c *Connection, vHost string) error {
	c.vHost = vHost
	return nil
}

// SSL toggles SSL/TLS on for this consumer's connection
func SSL(c *Connection, ssl bool) error {
	c.ssl = ssl
	return nil
}

// New - returns a new instance of the RabbitMQ provider
func New(options ...func(*Connection) error) (*Connection, error) {
	var err error
	var scheme string

	c := Connection{}

	// Default config
	c.host = "localhost"
	c.port = 5672
	c.user = "guest"
	c.password = "guest"

	// Process configuration options to override defaults
	for _, option := range options {
		// TODO - don't process NotifyClose until later
		err = option(&c)
		if err != nil {
			return nil, fmt.Errorf("failed to initialise connection: %s", err)
		}
	}

	if c.ssl {
		scheme = "amqps"
	} else {
		scheme = "amqp"
	}

	c.conn, err = amqp.Dial(
		fmt.Sprintf("%s://%s:%s@%s:%d/%s",
			scheme,
			c.user,
			c.password,
			c.host,
			c.port,
			c.vHost,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialise connection: %s", err)
	}

	return &c, nil
}

// Create a RabbitMQ channel from the connection
func (c *Connection) Channel() (*amqp.Channel, error) {
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("unable to create AMQP channel: %s", err)
	}

	return ch, nil
}

// Close close the connection to the RabbitMQ server
func (c *Connection) Close() error {
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("unable to close AMQP connection: %s", err)
	}

	return nil
}

//// NotifyClose registers a channel that will receive errors when
//// the connection to the RabbitMQ server is closed by the server
//func NotifyClose(ec chan error) func(c *Connection) error {
//	return func(c *Connection) error {
//		ac := make(chan *amqp.Error)
//
//		go func(ec chan error, ac chan *amqp.Error) {
//			// Pipe converted errors to our listener
//			for e := range ac {
//				ec <- error(e)
//			}
//
//			// When the AMQP error channel closes we need to close our channel
//			close(ec)
//		}(ec, ac)
//
//		return nil
//	}
//}
