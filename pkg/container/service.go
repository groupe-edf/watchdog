package container

// Service represents a service in the services container
type Service interface {
}

// ServiceProvider should be implemented by service providers, or such components, which register a service in the service container.
type ServiceProvider interface {
	Register(container Container)
}

// DefaultService default service implementation
type DefaultService struct {
	Shared bool
}

// IsShared check whether the service is shared or not
func (service *DefaultService) IsShared() bool {
	return service.Shared
}
