package commons

import (
	"net/rpc"
	"sync"
	"time"
)

// Wrapper struct that holds all the results of the checks for each account
type Tests struct {
	Account string  `yaml:"account"` // Account name
	Checks  []Check `yaml:"checks"`  // Checks
}

// A check is a test that is run on a resource
type Check struct {
	Name        string        `yaml:"name"`        // Name of the check
	Description string        `yaml:"description"` // Description of the check
	Status      string        `yaml:"status"`      // Status of the check - OK, FAIL
	Id          string        `yaml:"id"`          // ID of the check - unique identifier for the check - AWS_IAM_001
	Categories  []string      `yaml:"categories"`  // Category of the check - Security, Cost, Performance, Fault Tolerance, Operational Excellence, etc ...
	Results     []Result      `yaml:"results"`     // Results of the check
	Duration    time.Duration `yaml:"duration"`    // Duration of the check
	StartTime   time.Time
	EndTime     time.Time
}

type Result struct {
	Message    string `yaml:"message"`      // Message to display
	Status     string `yaml:"status"`       // Status of the check
	ResourceID string `yaml:"resource_arn"` // Resource ID - unique identifier for the resource
}

type Config struct {
	Plugins      []Plugin                 `yaml:"plugins"`
	Ignore       []Ignore                 `yaml:"ignore"`
	PluginConfig []map[string]interface{} `yaml:"pluginsConfiguration"`
	Tests        []Tests                  `yaml:"tests"`
}

// CheckConfig is a struct that contains all the information needed to run a check.
type CheckConfig struct {
	Wg          *sync.WaitGroup // Wait group to wait for all the checks to be done
	Queue       chan Check      // Queue to add the results to
	ConfigYatas *Config         // Yatas config
}

type Ignore struct {
	ID     string   `yaml:"id"`
	Regex  bool     `yaml:"regex"`
	Values []string `yaml:"values"`
}

// Yatas is the interface that we're exposing as a plugin.
type Yatas interface {
	Run(c *Config) []Tests
}

// Here is an implementation that talks over RPC
type YatasRPC struct{ client *rpc.Client }

// This is the implementation of plugin.Plugin so we can serve/consume this
//
// This has two methods: Server must return an RPC server for this plugin
// type. We construct a YatasRPCServer for this.
//
// Client must return an implementation of our interface that communicates
// over an RPC client. We return YatasRPC for this.
//
// Ignore MuxBroker. That is used to create more multiplexed streams on our
// plugin connection and is a more advanced use case.
type YatasPlugin struct {
	// Impl Injection
	Impl Yatas
}

// Here is the RPC server that YatasRPC talks to, conforming to
// the requirements of net/rpc
type YatasRPCServer struct {
	// This is the real implementation
	Impl Yatas
}

type Plugin struct {
	Name           string   `yaml:"name"`
	Enabled        bool     `yaml:"enabled"`
	Source         string   `yaml:"source"`
	Type           string   `default:"checks" yaml:"type" `
	Version        string   `yaml:"version"`
	Description    string   `yaml:"description"`
	Exclude        []string `yaml:"exclude"`
	Include        []string `yaml:"include"`
	Command        string   `yaml:"command"`
	Args           []string `yaml:"args"`
	ExpectedOutput string   `yaml:"expected_output"`
	ExpectedStatus int      `yaml:"expected_status"`

	// Parsed source attributes
	SourceOwner string
	SourceRepo  string
}
