package commons

import (
	"time"
)

type T Check

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
