package results

type T Check
type Result struct {
	Message    string `yaml:"message"`
	Status     string `yaml:"status"`
	ResourceID string `yaml:"resource_arn"`
}

type Check struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Status      string   `yaml:"status"`
	Id          string   `yaml:"id"`
	Results     []Result `yaml:"results"`
}

type Tests struct {
	Account string  `yaml:"account"`
	Checks  []Check `yaml:"checks"`
}
