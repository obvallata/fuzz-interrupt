package action

import (
	"fmt"
	"log"

	"diploma/chaosmachine/internal/interaction"
)

type breakpointHandler struct {
	action Action
	dlv    interaction.DlvClient
	config BreakpointConfig
}

type BreakpointConfig struct {
	Breakpoints []struct {
		FilePath string `yaml:"filePath"`
		Line     int    `yaml:"line"`
	} `yaml:"breakpoints"`
}

func (h *breakpointHandler) run() error {
	clearBreakpoints, err := h.createBreakpoints()
	if err != nil {
		return err
	}

	go func() {
		defer clearBreakpoints()

		for {
			if err != h.dlv.Continue() {
				log.Printf("Error dlv continue: %s\n", err.Error())
			}

			injectionName, err := h.dlv.GetVarStr("injectionName")
			if err != nil {
				log.Fatalf("Error dlv get injectionName: %s\n", err.Error())
			}

			h.action.HandleBreakpoint(injectionName)
		}
	}()

	return nil
}

func (h *breakpointHandler) createBreakpoints() (func(), error) {
	var ids []int

	for _, bp := range h.config.Breakpoints {
		id, err := h.dlv.CreateBreakpoint(bp.FilePath, bp.Line)
		if err != nil {
			return nil, fmt.Errorf("create breakpoint: %w", err)
		}
		ids = append(ids, id)
	}

	return func() {
		for _, id := range ids {
			// ignore error
			h.dlv.ClearBreakpoint(id)
		}
	}, nil
}
