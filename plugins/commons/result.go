package commons

type Result struct {
	Message    string `yaml:"message"`      // Message to display
	Status     string `yaml:"status"`       // Status of the check
	ResourceID string `yaml:"resource_arn"` // Resource ID - unique identifier for the resource
}
