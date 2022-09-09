package commons

import "time"

// A check is a test that is run on a resource
type Check struct {
	Name        string        `yaml:"name"`        // Name of the check
	Description string        `yaml:"description"` // Description of the check
	Status      string        `yaml:"status"`      // Status of the check - OK, FAIL
	Id          string        `yaml:"id"`          // ID of the check - unique identifier for the check - AWS_IAM_001
	Results     []Result      `yaml:"results"`     // Results of the check
	Duration    time.Duration `yaml:"duration"`    // Duration of the check
	StartTime   time.Time
	EndTime     time.Time
}
