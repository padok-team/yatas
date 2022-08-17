package report

import (
	"io/ioutil"
	"os"
	"strings"

	"testing"

	"github.com/stangirard/yatas/internal/results"
	"gopkg.in/yaml.v3"
)

func TestParseReportYaml(t *testing.T) {
	var report []results.Tests

	data, err := os.ReadFile("../testdata/results_data.yaml")
	if err != nil {
		t.Error(err)
	}
	err = yaml.Unmarshal(data, &report)
	if err != nil {
		t.Error(err)
	}
	if len(report) != 1 {
		t.Error("Expected 1 test, got", len(report))
	}

	if len(report[0].Checks) != 47 {
		t.Error("Expected 47 check, got", len(report[0].Checks))
	}

}

func TestGenerateReadme(t *testing.T) {
	// Catch the printed output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := GenerateReadme("../testdata/results_data.yaml")
	if err != nil {
		t.Error(err)
	}

	w.Close()
	os.Stdout = old
	out, _ := ioutil.ReadAll(r)
	//Read the file  readme_generated.txt and compare it to the output
	data, err := ioutil.ReadFile("../testdata/readme_generated.txt")
	if err != nil {
		t.Error(err)
	}
	// Replace all \n and space with nothing in data and out variables
	data = []byte(strings.Replace(string(data), "\n", "", -1))
	data = []byte(strings.Replace(string(data), " ", "", -1))
	out = []byte(strings.Replace(string(out), "\n", "", -1))
	out = []byte(strings.Replace(string(out), " ", "", -1))
	if string(out) != string(data) {
		t.Error("Expected:\n", string(data), "\nGot:\n", string(out))
	}

}
