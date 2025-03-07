package client

import (
	"diploma/keypoint/injection"
	"fmt"
	"github.com/monaco-io/request"
)

func (c *client) Enable(keyPointName string, config injection.Config) error {
	return c.send(&request.Client{
		URL:    c.url(keyPointName),
		Method: "PUT",
		JSON:   config,
	})
}

func (c *client) Disable(keyPointName string) error {
	return c.send(&request.Client{
		URL:    c.url(keyPointName),
		Method: "DELETE",
	})
}

func (c *client) send(httpClient *request.Client) error {
	resp := httpClient.Send()
	if !resp.OK() {
		return resp.Error()
	}
	return nil
}

func (c *client) url(keyPointName string) string {
	return fmt.Sprintf("%s/%s", c.config.URL, keyPointName)
}
