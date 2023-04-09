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

// RunPlugin runs the specified plugin with the given configuration and returns the test results.
func RunPlugin(pluginInput commons.Plugin, c *commons.Config) []commons.Tests {
	// Register types used for RPC communication
	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})

	// Create a plugin map to store available plugins
	pluginMap := make(map[string]plugin.Plugin)

	// Add plugins from the configuration to the plugin map
	for _, plugin := range c.Plugins {
		pluginMap[strings.ToLower(plugin.Name)] = &commons.YatasPlugin{}
	}

	// Expand the home directory for the plugin path
	homeDir, err := homedir.Expand("~/.yatas.d/plugins/")
	if err != nil {
		logger.Error("Error expanding home directory", "error", err)
	}
	// Construct the command to execute the plugin
	cmd := exec.Command(fmt.Sprintf("%s/%s/%s/yatas-%s", homeDir, pluginInput.Source, pluginInput.Version, pluginInput.Name))

	// Create a new plugin client with the specified configuration
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             cmd,
		Logger:          logger.Logger(),
	})
	defer client.Kill()

	// Connect to the plugin via RPC
	rpcClient, err := client.Client()
	if err != nil {
		if strings.Contains(err.Error(), "Incompatible API version with plugin") {
			logger.Error("Plugin is not compatible with YATAS. Please update it.", "plugin_name", pluginInput.Name, "error", err)
			return nil
		}
		logger.Error("Error creating RPC client", "error", err)
		return nil
	}

	// Request the plugin instance
	raw, err := rpcClient.Dispense(pluginInput.Name)
	if err != nil {
		logger.Error("Error dispensing plugin", "plugin_name", pluginInput.Name, "error", err)
		return nil
	}

	// Cast the received instance to the Yatas interface
	yatasPlugin := raw.(commons.Yatas)

	// Run the plugin and return the test results
	return yatasPlugin.Run(c)
}

// handshakeConfig is used for a basic handshake between the plugin and the host.
// This helps to prevent users from executing bad plugins or plugin directories.
// It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  2,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}
