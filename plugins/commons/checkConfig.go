package commons

import (
	"sync"
	"time"
)

// Init the check config struct. Particularly useful in the categories. It allows to pass the config to the checks and allows
// them to be run in parallel by adding the results to the queue.
func (c *CheckConfig) Init(config *Config) {
	c.Wg = &sync.WaitGroup{}
	c.Queue = make(chan Check, 10)
	c.ConfigYatas = config
}

type T Check

// Add Result to a check with some logic to update the status of the check
func (c *Check) AddResult(result Result) {
	if result.Status == "FAIL" {
		c.Status = "FAIL"
	}
	c.Results = append(c.Results, result)
}

// Initialise a check
func (c *Check) InitCheck(name, description, id string, categories []string) {
	c.Name = name
	c.Description = description
	c.Status = "OK"
	c.Id = id
	c.Categories = categories
	c.StartTime = time.Now()
}

// End a check by updating the duration and end time
func (c *Check) EndCheck() {
	c.EndTime = time.Now()
	c.Duration = c.EndTime.Sub(c.StartTime)
}
