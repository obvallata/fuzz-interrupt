package breakpoint

import (
	"fmt"
	"os"

	"diploma/keypoint/utils/ptr"
)

type Rollback func() error

func HandleInjection(config Config, data map[string]any) (Rollback, error) {
	switch config.Command {
	case FilePermissionChangeType:
		return handleFilePermissionChange(ptr.From(config.FilePermissionChangeConfig), data)
	default:
		return nil, fmt.Errorf("unexpected command")
	}
}

func handleFilePermissionChange(config FilePermissionChangeConfig, data map[string]any) (Rollback, error) {
	fmt.Println(config, data)

	filePath := config.FilePath
	if filePath == "" {
		filePath, _ = data[config.FilePathDataKey].(string)
	}

	if filePath == "" {
		return nil, fmt.Errorf("file path was not specified")
	}

	stat, err := os.Stat(filePath)
	fmt.Println(stat)
	if err != nil {
		return nil, err
	}
	prevFileMode := stat.Mode()

	if err := os.Chmod(filePath, config.FileMode); err != nil {
		return nil, err
	}

	fmt.Println("qweqwdqwdqw")

	return func() error {
		return os.Chmod(filePath, prevFileMode)
	}, nil
}
