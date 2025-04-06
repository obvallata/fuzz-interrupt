package interaction

import (
	"diploma/chaosmachine/internal/config"
	"diploma/keypoint/client"
)

type Clients struct {
	Dlv      DlvClient
	KeyPoint client.KeyPointClient
}

func NewClients(config config.InteractionConfig) (Clients, error) {
	return Clients{
		Dlv:      NewDlvClient(config.Dlv),
		KeyPoint: client.NewKeyPointClient(config.KeyPoint),
	}, nil
}

func (c *Clients) Close() error {
	return c.Dlv.Disconnect()
}
