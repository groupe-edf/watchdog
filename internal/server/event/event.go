package event

type Event interface {
	Abort(abort bool)
	Data() map[string]interface{}
	IsAborted() bool
	Name() string
	SetData(data map[string]interface{}) Event
	SetName(name string) Event
}

type EventType int

type DefaultEvent struct {
	data    map[string]interface{}
	name    string
	aborted bool
}

func (event *DefaultEvent) Abort(abort bool) {
	event.aborted = abort
}

func (event *DefaultEvent) Data() map[string]interface{} {
	return event.data
}

func (event *DefaultEvent) IsAborted() bool {
	return event.aborted
}

func (event *DefaultEvent) Name() string {
	return event.name
}

func (event *DefaultEvent) SetData(data map[string]interface{}) Event {
	if data != nil {
		event.data = data
	}
	return event
}

func (event *DefaultEvent) SetName(name string) Event {
	event.name = name
	return event
}
