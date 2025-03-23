package injection

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

type mockedStruct struct {
	Int    int         `json:"int"`
	Str    string      `json:"str"`
	Struct innerStruct `json:"struct"`
}

type innerStruct struct {
	Bool bool `json:"bool"`
}

func TestMock(t *testing.T) {
	f := reflect.TypeOf(func(a int) (*mockedStruct, mockedStruct, int, error) {
		return nil, mockedStruct{}, 0, nil
	})

	expectedPtr := &mockedStruct{
		Int:    42,
		Str:    "42",
		Struct: innerStruct{Bool: true},
	}
	expectedStruct := mockedStruct{Int: 69}
	expectedInt := 228
	expectedError := fmt.Errorf("privet")

	config := MockInjectionConfig{Outs: []any{expectedPtr, expectedStruct, expectedInt, expectedError}}

	values, err := Mock(config, f)
	require.NoError(t, err)
	require.Len(t, values, 4)

	vPtr, okPtr := values[0].Interface().(*mockedStruct)
	require.True(t, okPtr)
	require.Equal(t, expectedPtr, vPtr)

	vStruct, okStruct := values[1].Interface().(mockedStruct)
	require.True(t, okStruct)
	require.Equal(t, expectedStruct, vStruct)

	vInt, okInt := values[2].Interface().(int)
	require.True(t, okInt)
	require.Equal(t, expectedInt, vInt)

	vError, okError := values[3].Interface().(error)
	require.True(t, okError)
	require.Equal(t, expectedError, vError)
}
