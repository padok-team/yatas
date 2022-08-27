package yatas

import "time"

type T Check
type Result struct {
	Message    string `yaml:"message"`
	Status     string `yaml:"status"`
	ResourceID string `yaml:"resource_arn"`
}

type Check struct {
	Name        string        `yaml:"name"`
	Description string        `yaml:"description"`
	Status      string        `yaml:"status"`
	Id          string        `yaml:"id"`
	Results     []Result      `yaml:"results"`
	Duration    time.Duration `yaml:"duration"`
}

type Tests struct {
	Account string  `yaml:"account"`
	Checks  []Check `yaml:"checks"`
}

func (c *Check) AddResult(result Result) {
	if result.Status == "FAIL" {
		c.Status = "FAIL"
	}
	c.Results = append(c.Results, result)
}

func (c *Check) InitCheck(name, description, id string) {
	c.Name = name
	c.Description = description
	c.Status = "OK"
	c.Id = id
}
