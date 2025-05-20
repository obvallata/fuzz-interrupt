package breakpoint

import "os"

type Config struct {
	Command                    Command                     `yaml:"command"`
	FilePermissionChangeConfig *FilePermissionChangeConfig `yaml:"filePermissionChangeConfig"`
}

type Command string

const (
	FilePermissionChangeType Command = "file_permission_change"
)

type FilePermissionChangeConfig struct {
	// specify file path with one of the following ways:
	FilePath        string `yaml:"filePath"`        // static
	FilePathDataKey string `yaml:"filePathDataKey"` // will find file path in data

	FileMode os.FileMode `yaml:"fileMode"`
}
