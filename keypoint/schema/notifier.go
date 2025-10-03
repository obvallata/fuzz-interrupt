package schema

type NotifyRequest struct {
	Type NotifyType     `json:"type"`
	Name string         `json:"name"`
	Data map[string]any `json:"data"`
}

type NotifyType string

const (
	NotifyStateType     NotifyType = "state"
	NotifyInjectionType NotifyType = "injection"
)
