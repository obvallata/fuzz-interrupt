package interaction

import (
	"context"

	"diploma/keypoint/schema"
	"github.com/monaco-io/request"
)

type Notifier interface {
	Notify(ctx context.Context, req schema.NotifyRequest) error
}

var _ Notifier = (*notifierClient)(nil)

type notifierClient struct {
	config NotifierConfig
}

type NotifierConfig struct {
	URL string
}

func (c *notifierClient) Notify(ctx context.Context, req schema.NotifyRequest) error {
	cl := &request.Client{
		URL:    c.config.URL,
		Method: "POST",
		JSON:   req,
	}

	resp := cl.Send()
	if !resp.OK() {
		return resp.Error()
	}
	return nil
}

func NewNotifierClient(config NotifierConfig) *notifierClient {
	return &notifierClient{config: config}
}
