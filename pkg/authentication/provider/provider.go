package provider

// Provider authentication provider
type Provider interface {
	Authenticate(email string, password string) (*Identity, error)
}
