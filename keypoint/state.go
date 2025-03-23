package keypoint

import (
	"context"

	"diploma/keypoint/schema"
)

func State(ctx context.Context, stateName string, opts ...stateOption) {
	// TODO: add keypoint disabling env flag
	if !enabled.Load() {
		return
	}

	if notifier == nil {
		return
	}

	req := schema.NotifyRequest{
		Type: schema.NotifyStateType,
		Name: stateName,
	}
	for _, o := range opts {
		o(&req)
	}

	err := notifier.Notify(ctx, req)
	if err != nil {
		// TODO: handle error
	}
}

type stateOption func(req *schema.NotifyRequest)

func WithData(data map[string]any) stateOption {
	return func(req *schema.NotifyRequest) {
		req.Data = data
	}
}
