package injection

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
)

type MockInjectionConfig struct {
	Outs []any `json:"outs" yaml:"outs"`
}

// Mock supports only serializable types
func Mock(config MockInjectionConfig, f reflect.Type) ([]reflect.Value, error) {
	var outs []reflect.Value
	for i := 0; i < f.NumOut(); i++ {
		newP := reflect.New(f.Out(i)).Interface()
		if err := mapstructure.Decode(config.Outs[i], newP); err != nil {
			return nil, err
		}

		outs = append(outs, reflect.ValueOf(newP).Elem())
	}

	return outs, nil
}
