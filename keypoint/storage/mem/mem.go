package mem

import (
	"diploma/keypoint/injection"
	"diploma/keypoint/storage"
	"sync"
)

var _ storage.KeyPointStorage = (*Storage)(nil)

type Storage struct {
	mx sync.RWMutex
	m  map[string]injection.Config
}

func NewMemStorage() *Storage {
	return &Storage{m: make(map[string]injection.Config)}
}

func (s *Storage) GetInjectionConfig(keyPointName string) (injection.Config, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	return s.m[keyPointName], nil
}

func (s *Storage) UpdateInjectionConfig(keyPointName string, injectionConfig injection.Config) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.m[keyPointName] = injectionConfig
	return nil
}

func (s *Storage) Disable(keyPointName string) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	delete(s.m, keyPointName)
	return nil
}
