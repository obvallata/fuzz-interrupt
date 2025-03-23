package interaction

import "diploma/keypoint/client"

type Clients struct {
	Dlv      DlvClient
	KeyPoint client.KeyPointClient
}

type Config struct {
	Dlv      DlvConfig     `yaml:"dlv"`
	KeyPoint client.Config `yaml:"keyPoint"`
}

func NewClients(config Config) (Clients, error) {
	return Clients{
		Dlv:      NewDlvClient(config.Dlv),
		KeyPoint: client.NewKeyPointClient(config.KeyPoint),
	}, nil
}

func (c *Clients) Close() error {
	return c.Dlv.Disconnect()
}
