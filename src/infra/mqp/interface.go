package mqp

type MQPublisher interface {
	Publish(msg *Message) error
}

type ConsumerFunc func(msg *Message) error

func (sel ConsumerFunc) Process(msg *Message) error {
	return sel(msg)
}

type Consumer interface {
	ID() string
	Process(msg *Message) error
}
