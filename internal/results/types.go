package results

type T Check
type Result struct {
	Message    string `json:"message"`
	Status     string `json:"status"`
	ResourceID string `json:"resource_arn"`
}

type Check struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Status      string   `yaml:"status"`
	Id          string   `yaml:"id"`
	Results     []Result `yaml:"results"`
}

type Tests struct {
	Category string  `yaml:"category"`
	Checks   []Check `yaml:"checks"`
}
