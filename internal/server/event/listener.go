package event

type ListenerQueue struct {
	listeners []*ListenerItem
}

func (queue *ListenerQueue) Items() []*ListenerItem {
	return queue.listeners
}

func (queue *ListenerQueue) Push(listener *ListenerItem) *ListenerQueue {
	queue.listeners = append(queue.listeners, listener)
	return queue
}

type ListenerItem struct {
	Priority EventType
	Listener Listener
}

type Listener interface {
	Handle(event Event) error
}

type ListenerFunc func(event Event) error

func (callback ListenerFunc) Handle(event Event) error {
	return callback(event)
}
