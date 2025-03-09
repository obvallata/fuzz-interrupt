package storage

import (
	"diploma/keypoint/injection"
)

type KeyPointStorage interface {
	GetInjectionConfig(keyPointName string) (injection.Config, error)
	UpdateInjectionConfig(keyPointName string, injectionConfig injection.Config) error
	Disable(keyPointName string) error
}
