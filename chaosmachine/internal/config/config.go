package config

import (
	"fmt"
	"os"

	"diploma/chaosmachine/internal/breakpoint"
	"diploma/keypoint/client"
	"diploma/keypoint/injection"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server     ServerConfig      `yaml:"server"`
	Clients    InteractionConfig `yaml:"clients"`
	Breakpoint BreakpointConfig  `yaml:"breakpoint"`

	Injections InjectionsConfig       `yaml:"injections"`
	States     map[string]StateConfig `yaml:"states"`
}

type InjectionsConfig map[string]InjectionConfig

type InjectionConfig struct {
	InjectionsList []InjectionListElement `yaml:"injectionsList"`
}

type BreakpointInjectionConfig = breakpoint.Config

type InjectionListElement struct {
	injection.Config          `yaml:",inline"`
	BreakpointInjectionConfig `yaml:",inline"`
}

type StateConfig struct {
	Finish bool `yaml:"finish"`
}

type ServerConfig struct {
	URL string `yaml:"url"`
}

type InteractionConfig struct {
	Dlv      InteractionDlvConfig `yaml:"dlv"`
	KeyPoint client.Config        `yaml:"keyPoint"`
}

type InteractionDlvConfig struct {
	Host string `yaml:"host"`
}

type BreakpointConfig struct {
	Breakpoints []struct {
		FilePath string `yaml:"filePath"`
		Line     int    `yaml:"line"`
	} `yaml:"breakpoints"`
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
