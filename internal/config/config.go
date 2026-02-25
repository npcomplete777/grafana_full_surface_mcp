// Package config handles tool enable/disable configuration for the Grafana MCP server.
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ToolConfig controls whether an individual MCP tool is enabled.
type ToolConfig struct {
	Enabled *bool `yaml:"enabled"`
}

// yamlConfig is the raw YAML file structure.
type yamlConfig struct {
	Tools map[string]ToolConfig `yaml:"tools"`
}

// ToolsConfig holds per-tool enable/disable settings loaded from a YAML file.
type ToolsConfig struct {
	tools map[string]ToolConfig
}

// Load reads tool configuration from the file pointed to by GRAFANA_CONFIG_FILE,
// falling back to config.yaml in the working directory. A missing file is silently
// ignored — all tools default to enabled.
func Load() (*ToolsConfig, error) {
	path := os.Getenv("GRAFANA_CONFIG_FILE")
	if path == "" {
		path = "config.yaml"
	}

	cfg := &ToolsConfig{tools: make(map[string]ToolConfig)}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var y yamlConfig
	if err := yaml.Unmarshal(data, &y); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}
	if y.Tools != nil {
		cfg.tools = y.Tools
	}
	return cfg, nil
}

// IsEnabled reports whether the named tool should be registered.
// Tools absent from the config file default to enabled.
func (c *ToolsConfig) IsEnabled(name string) bool {
	tc, ok := c.tools[name]
	if !ok {
		return true
	}
	if tc.Enabled == nil {
		return true
	}
	return *tc.Enabled
}
