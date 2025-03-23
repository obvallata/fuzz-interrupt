package client

import (
	"diploma/keypoint/injection"
)

type KeyPointClient interface {
	Enable(keyPointName string, config injection.Config) error
	Disable(keyPointName string) error
}

var _ KeyPointClient = (*client)(nil)

type client struct {
	config Config
}

type Config struct {
	URL string
}

func NewKeyPointClient(config Config) *client {
	return &client{config: config}
}
