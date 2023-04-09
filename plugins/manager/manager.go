package manager

import (
	"encoding/gob"
	"fmt"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-plugin"
	"github.com/mitchellh/go-homedir"
	"github.com/padok-team/yatas/plugins/commons"
	"github.com/padok-team/yatas/plugins/logger"
)

func RunPlugin(pluginInput commons.Plugin, c *commons.Config) []commons.Tests {
	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})
	pluginMap := make(map[string]plugin.Plugin)

	for _, plugin := range c.Plugins {
		pluginMap[strings.ToLower(plugin.Name)] = &commons.YatasPlugin{}
	}

	homeDir, err := homedir.Expand("~/.yatas.d/plugins/")
	if err != nil {
		logger.Error("Error expanding home directory", "error", err)
	}
	cmd := exec.Command(fmt.Sprintf("%s/%s/%s/yatas-%s", homeDir, pluginInput.Source, pluginInput.Version, pluginInput.Name))

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             cmd,
		Logger:          logger.Logger(),
	})
	defer client.Kill()

	rpcClient, err := client.Client()
	if err != nil {
		if strings.Contains(err.Error(), "Incompatible API version with plugin") {
			logger.Error("Plugin is not compatible with YATAS. Please update it.", "plugin_name", pluginInput.Name, "error", err)
			return nil
		}
		logger.Error("Error creating RPC client", "error", err)
		return nil
	}

	raw, err := rpcClient.Dispense(pluginInput.Name)
	if err != nil {
		logger.Error("Error dispensing plugin", "plugin_name", pluginInput.Name, "error", err)
		return nil
	}

	yatasPlugin := raw.(commons.Yatas)

	return yatasPlugin.Run(c)
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  2,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}
