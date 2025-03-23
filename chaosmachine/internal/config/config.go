package config

import (
	"fmt"
	"os"

	"diploma/chaosmachine/internal/action"
	"diploma/chaosmachine/internal/interaction"
	"diploma/chaosmachine/internal/server"
	"diploma/keypoint/injection"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server     server.Config           `yaml:"server"`
	Clients    interaction.Config      `yaml:"clients"`
	Breakpoint action.BreakpointConfig `yaml:"breakpoint"`

	Injections map[string][]injection.Config `yaml:"injections"`
}

func GetConfig(configPath string) (Config, error) {
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("read config file: %w", err)
	}

	var c Config
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return Config{}, fmt.Errorf("unmarshal config: %w", err)
	}

	return c, nil
}
