package caller

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
)

const (
	magicFileWithProcessPID = "/Users/ddr/fuzz-interrupt/agent/src/pid"
)

type SignalCaller struct {
}

func NewSignalCaller() *SignalCaller {
	return &SignalCaller{}
}

func (c *SignalCaller) Call() error {
	pid, err := c.getPID()
	if err != nil {
		return fmt.Errorf("failed to get pid: %w", err)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process: %w", err)
	}

	err = process.Signal(syscall.SIGINFO)
	if err != nil {
		return fmt.Errorf("failed to signal process: %w", err)
	}

	return nil
}

func (c *SignalCaller) getPID() (int, error) {
	pidByte, err := os.ReadFile(magicFileWithProcessPID)
	if err != nil {
		return -1, fmt.Errorf("failed to read file: %w", err)
	}
	pidStr := strings.TrimSpace(string(pidByte))
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return -1, fmt.Errorf("failed to convert pid from string: %w", err)
	}

	return pid, nil
}
