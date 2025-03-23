package keypoint

import (
	"context"
	"reflect"
	"time"

	"diploma/keypoint/injection"
	"diploma/keypoint/schema"
	"diploma/keypoint/utils/ptr"
)

func WithInject[T any](ctx context.Context, injectionName string, function T) T {
	// TODO: add keypoint disabling env flag
	if !enabled.Load() {
		return function
	}

	v := reflect.MakeFunc(reflect.TypeOf(function), func(in []reflect.Value) []reflect.Value {
		notifyInjection(ctx, injectionName)
		outs, err := callWithInject(ctx, injectionName, function, in)

		if err != nil {
			// TODO: handle error
		}
		return outs
	})

	return v.Interface().(T)
}

func callWithInject[T any](
	ctx context.Context,
	injectName string,
	function T,
	in []reflect.Value,
) ([]reflect.Value, error) {
	injectionConfig, err := keyPointStorage.GetInjectionConfig(injectName)
	if err != nil {
		return originalCall(function, in), err
	}

	// TODO: hide in injection package
	switch injectionConfig.Type {
	case injection.TypeSleep:
		time.Sleep(ptr.From(injectionConfig.Sleep).Duration)

	case injection.TypeMock:
		outs, err := injection.Mock(ptr.From(injectionConfig.Mock), reflect.TypeOf(function))
		if err != nil {
			return originalCall(function, in), err
		}
		return outs, nil

	case injection.TypeBreakpoint:
		injection.Breakpoint(injectName, ptr.From(injectionConfig.Breakpoint))

	default: // including injection.TypeOff
	}

	return originalCall(function, in), nil
}

func notifyInjection(ctx context.Context, injectionName string) {
	if notifier == nil {
		return
	}

	err := notifier.Notify(ctx, schema.NotifyRequest{
		Type: schema.NotifyInjectionType,
		Name: injectionName,
	})
	if err != nil {
		// TODO: handle error
	}
}

func originalCall[T any](function T, in []reflect.Value) []reflect.Value {
	f := reflect.ValueOf(function)
	return f.Call(in)
}
