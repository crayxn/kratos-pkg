package queue

type Message struct {
	Topic   string
	Payload []byte
}

type Producer interface {
	Send(Message) error
}

type Consumers interface {
	RegisterHandler(topic string, handler func(Message) error)
	Start() error
	Stop() error
}
