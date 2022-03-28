package event

import (
	"sync"

	"github.com/minio/pkg/wildcard"
)

type Event struct {
	Data  interface{}
	Topic string
	wg    *sync.WaitGroup
}

func (event *Event) Done() {
	if event.wg != nil {
		event.wg.Done()
	}
}

type EventChannel chan Event

type eventChannelSlice []EventChannel

type CallbackFunc func(topic string, data interface{})

type EventBus struct {
	lock        sync.RWMutex
	subscribers map[string]eventChannelSlice
}

func (bus *EventBus) Publish(topic string, data interface{}) interface{} {
	wg := sync.WaitGroup{}
	channels := bus.getSubscribingChannels(topic)
	wg.Add(len(channels))
	bus.publish(channels, Event{
		Data:  data,
		Topic: topic,
		wg:    &wg,
	})
	wg.Wait()
	return data
}

func (bus *EventBus) PublishAsync(topic string, data interface{}) {
	bus.publish(bus.getSubscribingChannels(topic), Event{
		Data:  data,
		Topic: topic,
		wg:    nil,
	})
}

func (bus *EventBus) Subscribe(topic string) EventChannel {
	channel := make(EventChannel)
	bus.SubscribeChannel(topic, channel)
	return channel
}

func (bus *EventBus) SubscribeChannel(topic string, channel EventChannel) {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	if previous, found := bus.subscribers[topic]; found {
		bus.subscribers[topic] = append(previous, channel)
	} else {
		bus.subscribers[topic] = append([]EventChannel{}, channel)
	}
}

func (bus *EventBus) SubscribeCallback(topic string, callable CallbackFunc) {
	channel := NewEventChannel()
	bus.SubscribeChannel(topic, channel)
	go func(callable CallbackFunc) {
		for {
			event := <-channel
			callable(event.Topic, event.Data)
			event.Done()
		}
	}(callable)
}

func (bus *EventBus) publish(channels eventChannelSlice, event Event) {
	bus.lock.RLock()
	defer bus.lock.RUnlock()
	go func(channels eventChannelSlice, event Event) {
		for _, channel := range channels {
			channel <- event
		}
	}(channels, event)
}

func (eb *EventBus) getSubscribingChannels(topic string) eventChannelSlice {
	subChannels := eventChannelSlice{}
	for topicName := range eb.subscribers {
		if topicName == topic || wildcard.Match(topicName, topic) {
			subChannels = append(subChannels, eb.subscribers[topicName]...)
		}
	}
	return subChannels
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: map[string]eventChannelSlice{},
	}
}

func NewEventChannel() EventChannel {
	return make(EventChannel)
}
