package testreporter

import (
	"fmt"
	"go.uber.org/mock/gomock"
)

// Он существует, так как моку нужен TestReporter, а из рамок go test мы вышли.

var _ gomock.TestReporter = (*TestReporter)(nil)

type TestReporter struct{}

func (t *TestReporter) Errorf(format string, args ...any) {
	fmt.Printf(format, args...)
}

func (t *TestReporter) Fatalf(format string, args ...any) {
	fmt.Printf(format, args...)
}
