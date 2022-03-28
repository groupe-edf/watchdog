package container

import (
	"fmt"
	"sync"
)

var di Container

// ContainerFunc type
type ContainerFunc func(container Container) Service

// Container container interface
type Container interface {
	Get(service string) Service
	Has(service string) bool
	Provide(provider ServiceProvider)
	Set(service string, function ContainerFunc)
}

// ServiceContainer service container
type ServiceContainer struct {
	definitions map[string]ContainerFunc
	lock        *sync.RWMutex
	services    map[string]Service
}

// Get return service
func (container *ServiceContainer) Get(service string) Service {
	container.lock.RLock()
	_, ok := container.definitions[service]
	container.lock.RUnlock()
	if !ok {
		panic(fmt.Sprintf("The service does not exist: %s", service))
	}
	container.lock.RLock()
	_, ok = container.services[service]
	container.lock.RUnlock()
	if !ok {
		container.lock.Lock()
		container.services[service] = container.definitions[service](container)
		container.lock.Unlock()
	}
	container.lock.RLock()
	defer container.lock.RUnlock()
	return container.services[service]
}

// Has check if service exists in container
func (container *ServiceContainer) Has(service string) bool {
	container.lock.RLock()
	defer container.lock.RUnlock()
	return false
}

// Provide add service provider
func (container *ServiceContainer) Provide(provider ServiceProvider) {
	provider.Register(container)
}

// Set register a service
func (container *ServiceContainer) Set(service string, function ContainerFunc) {
	container.lock.Lock()
	defer container.lock.Unlock()
	container.definitions[service] = function
}

// Provide add service provider
func Provide(provider ServiceProvider) {
	di.Provide(provider)
}

// Get return service
func Get(service string) Service {
	return di.Get(service)
}

// GetContainer return registered container
func GetContainer() Container {
	return di
}

// NewServiceContainer retrun new services container
func NewServiceContainer() Container {
	return &ServiceContainer{
		definitions: make(map[string]ContainerFunc),
		lock:        &sync.RWMutex{},
		services:    make(map[string]Service),
	}
}

func init() {
	di = NewServiceContainer()
}
