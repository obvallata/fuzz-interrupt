package client

import (
	"diploma/keypoint/injection"
	"diploma/keypoint/schema"
)

type KeyPointClient interface {
	EnableMonitor(request schema.EnableMonitorRequest) error
	DisableMonitor() error

	EnableInjection(keyPointName string, config injection.Config) error
	DisableInjection(keyPointName string) error
}

var _ KeyPointClient = (*client)(nil)

type client struct {
	config Config
}

type Config struct {
	URL string `yaml:"url"`
}

func NewKeyPointClient(config Config) *client {
	return &client{config: config}
}
