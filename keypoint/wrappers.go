package keypoint

import (
	"context"
	"diploma/keypoint/injection"
	"diploma/keypoint/utils/ptr"
	"reflect"
	"time"
)

func WithInject[T any](ctx context.Context, keypointName string, function T) T {
	// TODO: keypoint disabling env flag

	v := reflect.MakeFunc(reflect.TypeOf(function), func(in []reflect.Value) []reflect.Value {
		notifyStart(keypointName)
		outs, err := callWithInject(ctx, keypointName, function, in)

		if err != nil {
			notifyError(keypointName)
		} else {
			notifySuccess(keypointName)
		}

		return outs
	})

	return v.Interface().(T)
}

func callWithInject[T any](
	ctx context.Context,
	keypointName string,
	function T,
	in []reflect.Value,
) ([]reflect.Value, error) {
	injectionConfig, err := keyPointStorage.GetInjectionConfig(keypointName)
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
		injection.Breakpoint(keypointName, ptr.From(injectionConfig.Breakpoint))

	default: // including injection.TypeOff
	}

	return originalCall(function, in), nil
}

func notifyStart(keypointName string) {
	injection.Breakpoint(keypointName, injection.BreakpointInjectionConfig{Command: injection.BreakpointNotifyStartType})
}

func notifySuccess(keypointName string) {
	injection.Breakpoint(keypointName, injection.BreakpointInjectionConfig{Command: injection.BreakpointNotifySuccessType})
}

func notifyError(keypointName string) {
	injection.Breakpoint(keypointName, injection.BreakpointInjectionConfig{Command: injection.BreakpointNotifyErrorType})
}

func originalCall[T any](function T, in []reflect.Value) []reflect.Value {
	f := reflect.ValueOf(function)
	return f.Call(in)
}
