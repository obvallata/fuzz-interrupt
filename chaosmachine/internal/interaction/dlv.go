package interaction

import (
	"fmt"

	"github.com/go-delve/delve/service/api"
	"github.com/go-delve/delve/service/rpc2"
)

type DlvClient interface {
	// CreateBreakpoint creates breakpoint and return its ID
	CreateBreakpoint(filePath string, line int) (int, error)
	ClearBreakpoint(id int) error

	// Continue synchronously waits when program achieve breakpoint
	Continue() error

	// GetVarStr inspects string variable value
	GetVarStr(varName string) (string, error)

	// Disconnect should be always called for closing interaction.
	// It also calls `continue` to not block on breakpoint
	Disconnect() error
}

var _ DlvClient = (*dlvClient)(nil)

type dlvClient struct {
	*rpc2.RPCClient
}

type DlvConfig struct {
	Host string `yaml:"host"`
}

func NewDlvClient(config DlvConfig) DlvClient {
	return &dlvClient{rpc2.NewClient(config.Host)}
}

var defaultLoadConfig = api.LoadConfig{
	FollowPointers:     true,
	MaxVariableRecurse: 1,
	MaxStringLen:       64,
	MaxArrayValues:     64,
	MaxStructFields:    -1,
}

func (c *dlvClient) CreateBreakpoint(filePath string, line int) (int, error) {
	// Stop program
	_, err := c.Halt()
	if err != nil {
		return 0, fmt.Errorf("dlv halt: %w", err)
	}

	bp, err := c.RPCClient.CreateBreakpoint(&api.Breakpoint{
		File:       filePath,
		Line:       line,
		Tracepoint: true,
	})
	if err != nil {
		return 0, fmt.Errorf("dlv create breakpoint: %w", err)
	}

	// Continue program
	c.RPCClient.Continue()

	return bp.ID, nil
}

func (c *dlvClient) ClearBreakpoint(id int) error {
	_, err := c.RPCClient.ClearBreakpoint(id)
	return err
}

func (c *dlvClient) Continue() error {
	out := new(rpc2.CommandOut)
	return c.CallAPI("Command", &api.DebuggerCommand{Name: api.Continue, ReturnInfoLoadConfig: &defaultLoadConfig}, &out)
}

func (c *dlvClient) GetVarStr(varName string) (string, error) {
	varCommand, err := c.EvalVariable(api.EvalScope{GoroutineID: -1}, varName, defaultLoadConfig)
	if err != nil {
		return "", fmt.Errorf("dlv eval variable: %w", err)
	}

	return varCommand.Value, nil
}

func (c *dlvClient) Disconnect() error {
	return c.RPCClient.Disconnect(true)
}
