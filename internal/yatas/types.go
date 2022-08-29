package yatas

import (
	"time"
)

type T Check

// Result is a specific result of a check for a given resource
type Result struct {
	Message    string `yaml:"message"`      // Message to display
	Status     string `yaml:"status"`       // Status of the check
	ResourceID string `yaml:"resource_arn"` // Resource ID - unique identifier for the resource
}

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

// Wrapper struct that holds all the results of the checks for each account
type Tests struct {
	Account string  `yaml:"account"` // Account name
	Checks  []Check `yaml:"checks"`  // Checks
}

// Add Result to a check with some logic to update the status of the check
func (c *Check) AddResult(result Result) {
	if result.Status == "FAIL" {
		c.Status = "FAIL"
	}
	c.Results = append(c.Results, result)
}

// Initialise a check
func (c *Check) InitCheck(name, description, id string) {
	c.Name = name
	c.Description = description
	c.Status = "OK"
	c.Id = id
	c.StartTime = time.Now()
}

// End a check by updating the duration and end time
func (c *Check) EndCheck() {
	c.EndTime = time.Now()
	c.Duration = c.EndTime.Sub(c.StartTime)
}
