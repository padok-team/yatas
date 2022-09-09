package commons

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
)

// CheckConfig is a struct that contains all the information needed to run a check.
type CheckConfig struct {
	Wg          *sync.WaitGroup // Wait group to wait for all the checks to be done
	ConfigAWS   aws.Config      // AWS config
	Queue       chan Check      // Queue to add the results to
	ConfigYatas *Config         // Yatas config
}

// Init the check config struct. Particularly useful in the categories. It allows to pass the config to the checks and allows
// them to be run in parallel by adding the results to the queue.
func (c *CheckConfig) Init(s aws.Config, config *Config) {
	c.Wg = &sync.WaitGroup{}
	c.ConfigAWS = s
	c.Queue = make(chan Check, 10)
	c.ConfigYatas = config
}
