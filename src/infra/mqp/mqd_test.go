package mqp

import (
	"fmt"
	"testing"
	"time"
)

type MockConsumer struct {
}

func (m *MockConsumer) ID() string {
	return "MockConsumer"
}

func (m *MockConsumer) Process(msg *Message) error {
	fmt.Println("收到消息", msg)
	return nil
}

func TestNewMDP(t *testing.T) {
	mdp := NewMDP(200, nil)
	mdp.SubScribe(&MockConsumer{})
	mdp.Start()
	mdp.Publish(&Message{})
	time.Sleep(time.Second * 1)

}
