package event

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testListener struct {
}

func (listener *testListener) Handle(event Event) error {
	return fmt.Errorf("unexpected error")
}

func TestManager(t *testing.T) {
	manager := NewManager()
	manager.On("analysis", &ListenerItem{
		Listener: &testListener{},
		Priority: EventType(100),
	})
	event := &DefaultEvent{}
	event.SetName("analysis")
	manager.Attach(event)
	_, err := manager.Fire("analysis", nil)
	assert.NoError(t, err)
}
