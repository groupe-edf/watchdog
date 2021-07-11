package event

type Manager interface {
	Attach(event Event)
	Detach(name string)
	On(name string, listener *ListenerItem, priority ...int)
	Fire(name string, data map[string]interface{}) (event Event, err error)
}

type EventData struct {
}

type HandlerFunc func(data *EventData) error

type EventManager struct {
	events    map[string]Event
	listeners map[string]*ListenerQueue
	names     map[string]int
}

func (manager *EventManager) AddListener(name string, listener *ListenerItem, priority ...int) {
	manager.On(name, listener, priority...)
}

func (manager *EventManager) Attach(event Event) {
	manager.events[event.Name()] = event
}

func (manager *EventManager) Detach(name string) {
	delete(manager.events, name)
}

func (manager *EventManager) Fire(name string, data map[string]interface{}) (event Event, err error) {
	if event, ok := manager.events[name]; ok {
		if data != nil {
			event.SetData(data)
			err = manager.FireEvent(event)
			return event, err
		}
	}
	return event, err
}

func (manager *EventManager) FireEvent(event Event) (err error) {
	event.Abort(false)
	queue, ok := manager.listeners[event.Name()]
	if ok {
		for _, listener := range queue.Items() {
			err = listener.Listener.Handle(event)
			if err != nil || event.IsAborted() {
				return
			}
		}
	}
	return
}

func (manager *EventManager) GetEvent(name string) (event Event, ok bool) {
	event, ok = manager.events[name]
	return
}

func (manager *EventManager) HasEvent(name string) bool {
	_, ok := manager.events[name]
	return ok
}

func (manager *EventManager) On(name string, listener *ListenerItem, priority ...int) {
	if queue, ok := manager.listeners[name]; ok {
		queue.Push(listener)
	} else {
		manager.names[name] = 1
		manager.listeners[name] = (&ListenerQueue{}).Push(listener)
	}
}

func NewManager() Manager {
	return &EventManager{
		events:    make(map[string]Event),
		listeners: make(map[string]*ListenerQueue),
		names:     make(map[string]int),
	}
}
