package custom

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/stangirard/yatas/plugins/commons"
)

func Run(c *commons.Config, name string) (commons.Tests, error) {
	plugin := c.FindPluginWithName(name)
	checks, err := ExecuteCommand(c, plugin)
	return checks, err

}

func ExecuteCommand(c *commons.Config, plugin *commons.Plugin) (commons.Tests, error) {
	checks := []commons.Check{}
	check := commons.Check{}
	check.Name = plugin.Name
	check.Description = plugin.Description
	check.Status = "OK"

	cmd := exec.Command(plugin.Command, plugin.Args...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	result := commons.Result{}
	if strings.TrimRight(outb.String(), "\n") == plugin.ExpectedOutput {
		result.Message = fmt.Sprint("Output matched: ", plugin.ExpectedOutput)
		result.Status = "OK"
	} else {
		result.Message = fmt.Sprint("Output did not match: ", plugin.ExpectedOutput, " instead got: ", outb.String())
		result.Status = "FAIL"
	}
	check.Results = append(check.Results, result)
	checks = append(checks, check)
	test := commons.Tests{}
	test.Checks = checks
	test.Account = plugin.Name
	return test, nil

}
