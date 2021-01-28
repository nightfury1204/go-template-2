package rabbitmq

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"log"
	"sync"
	"time"
)

type Rabbit struct {
	URL            string
	connectionName string
	conn           *amqp.Connection
	mux            sync.RWMutex
	Reconnect      struct {
		Interval   time.Duration
		MaxAttempt int
	}
}

func New(url, name string) *Rabbit {
	rbt := &Rabbit{
		URL:            url,
		connectionName: name,
	}
	rbt.Reconnect.Interval = 500 * time.Millisecond
	rbt.Reconnect.MaxAttempt = 7200
	return rbt
}

// Connect connects to RabbitMQ server.
func (r *Rabbit) Connect() error {
	if r.conn == nil || r.conn.IsClosed() {
		con, err := amqp.DialConfig(r.URL, amqp.Config{Properties: amqp.Table{"connection_name": r.connectionName}})
		if err != nil {
			return err
		}
		r.conn = con
	}
	//go r.reconnect()

	return nil
}

func (r *Rabbit) Ping() error {
	if r.conn.IsClosed() {
		return fmt.Errorf("connection_name: %s closed", r.connectionName)
	}
	return nil
}

func (r *Rabbit) Connection() (*amqp.Connection, error) {
	if r.conn == nil || r.conn.IsClosed() {
		return nil, errors.New("connection is not open")
	}

	return r.conn, nil
}

// Channel returns a new `*amqp.Channel` instance.
func (r *Rabbit) Channel() (*amqp.Channel, error) {
	chn, err := r.conn.Channel()
	if err != nil {
		return nil, err
	}

	return chn, nil
}

func (r *Rabbit) Shutdown() error {
	if r.conn != nil {
		return r.conn.Close()
	}

	return nil
}

// reconnect reconnects to server if the connection or a channel
// is closed unexpectedly. Normal shutdown is ignored. It tries
// maximum of 7200 times and sleeps half a second in between
// each try which equals to 1 hour.
func (r *Rabbit) reconnect() {
WATCH:

	conErr := <-r.conn.NotifyClose(make(chan *amqp.Error))
	if conErr != nil {
		log.Println("CRITICAL: Connection dropped, reconnecting")

		var err error

		for i := 1; i <= r.Reconnect.MaxAttempt; i++ {
			r.mux.RLock()
			r.conn, err = amqp.DialConfig(r.URL, amqp.Config{Properties: amqp.Table{"connection_name": r.connectionName}})
			r.mux.RUnlock()

			if err == nil {
				log.Println("INFO: Reconnected")

				goto WATCH
			}

			time.Sleep(r.Reconnect.Interval)
		}

		log.Println(errors.Wrap(err, "CRITICAL: Failed to reconnect"))
	} else {
		log.Println("INFO: Connection dropped normally, will not reconnect")
	}
}
