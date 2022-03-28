package hook

// GitHooks data structure
type GitHooks struct {
	Hooks   []Hook `yaml:"hooks,omitempty"`
	Locked  bool   `yaml:"locked"`
	Version string `yaml:"version,omitempty"`
}

// Hook hook aggregate model
type Hook struct {
	Description      string  `yaml:"description"`
	Name             string  `yaml:"name,omitempty"`
	RejectionMessage string  `yaml:"rejection_message"`
	Rules            []*Rule `yaml:"rules,omitempty"`
}
