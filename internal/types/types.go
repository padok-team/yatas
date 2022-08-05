package types

type Result struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type Check struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Status      string   `yaml:"status"`
	Results     []Result `yaml:"results"`
}

type Tests struct {
	Category string  `yaml:"category"`
	Checks   []Check `yaml:"checks"`
}
