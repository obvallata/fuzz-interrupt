package client

import (
	"fmt"

	"diploma/keypoint/injection"
	"diploma/keypoint/schema"
	"github.com/monaco-io/request"
)

func (c *client) EnableMonitor(req schema.EnableMonitorRequest) error {
	return c.send(&request.Client{
		URL:    fmt.Sprintf("%s/monitor/enable", c.config.URL),
		Method: "POST",
		JSON:   req,
	})
}

func (c *client) DisableMonitor() error {
	return c.send(&request.Client{
		URL:    fmt.Sprintf("%s/monitor/disable", c.config.URL),
		Method: "POST",
	})
}

func (c *client) EnableInjection(keyPointName string, config injection.Config) error {
	return c.send(&request.Client{
		URL:    c.injectionURL(keyPointName),
		Method: "PUT",
		JSON:   config,
	})
}

func (c *client) DisableInjection(keyPointName string) error {
	return c.send(&request.Client{
		URL:    c.injectionURL(keyPointName),
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

func (c *client) injectionURL(keyPointName string) string {
	return fmt.Sprintf("%s/injection/%s", c.config.URL, keyPointName)
}
