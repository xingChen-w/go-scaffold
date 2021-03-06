package producer

import (
	"sync"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
)

var (
	ErrProducerClosed = errors.New("[kafka] the producer is closed")
)

type (
	// ProducerHandlerError triggered every time when a error message is reveived.
	ProducerHandlerError func(*sarama.ProducerError)
)

// Producer ..
type AsyncProducer struct {
	config *Config
	p      sarama.AsyncProducer

	closed bool

	rwmu sync.RWMutex
}

// create async producer instance.
// create a new goroutine for watch the prodcer.Errors channel and
// triggered CallbackError handle every time when a error message is received.
func newAsyncProducer(conf *Config) (*AsyncProducer, error) {
	version, err := sarama.ParseKafkaVersion(conf.Version)
	if err != nil {
		return nil, err
	}
	c := sarama.NewConfig()
	c.Version = version
	c.Producer.RequiredAcks = sarama.WaitForAll
	c.Producer.Partitioner = sarama.NewRandomPartitioner
	producer, err := sarama.NewAsyncProducer(conf.Brokers, c)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			err, ok := <-producer.Errors()
			if !ok {
				return
			}
			if conf.CallbackError != nil {
				conf.CallbackError(err)
			}
		}
	}()

	return &AsyncProducer{
		config: conf,
		p:      producer,
		closed: false,
	}, nil
}

// Send message to topic.
// return ProducerErrClosed error when the producer is closed.
func (t *AsyncProducer) Send(message string) error {
	t.rwmu.RLock()
	defer t.rwmu.RUnlock()
	if t.closed {
		return ErrProducerClosed
	}
	t.p.Input() <- &sarama.ProducerMessage{
		Topic: t.config.Topic,
		Value: sarama.StringEncoder(message),
	}
	return nil
}

// Close async producer.
func (t *AsyncProducer) Close() error {
	t.rwmu.Lock()
	defer t.rwmu.Unlock()
	if t.closed {
		return nil
	}
	err := t.p.Close()
	if err != nil {
		return err
	}
	t.closed = true
	return nil
}

// Producer ..
type Producer struct {
	config *Config
	p      sarama.SyncProducer

	closed bool

	rwmu sync.RWMutex
}

// create sync producer instance.
func newProducer(conf *Config) (*Producer, error) {
	version, err := sarama.ParseKafkaVersion(conf.Version)
	if err != nil {
		return nil, err
	}
	c := sarama.NewConfig()
	c.Version = version
	c.Producer.RequiredAcks = sarama.WaitForAll
	c.Producer.Partitioner = sarama.NewRandomPartitioner
	c.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(conf.Brokers, c)
	if err != nil {
		return nil, err
	}
	return &Producer{
		p:      producer,
		config: conf,
		closed: false,
	}, nil
}

// Send message to topic.
// return ProducerErrClosed error when the producer is closed.
func (t *Producer) Send(message string) error {
	t.rwmu.RLock()
	defer t.rwmu.RUnlock()
	if t.closed {
		return ErrProducerClosed
	}
	_, _, err := t.p.SendMessage(&sarama.ProducerMessage{
		Topic: t.config.Topic,
		Value: sarama.StringEncoder(message),
	})
	return err
}

// Close producer.
func (t *Producer) Close() error {
	t.rwmu.Lock()
	defer t.rwmu.Unlock()
	if t.closed {
		return nil
	}
	err := t.p.Close()
	if err != nil {
		return err
	}
	t.closed = true
	return nil
}
