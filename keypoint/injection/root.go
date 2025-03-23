package injection

type Config struct {
	Type Type `json:"type"`

	Sleep      *SleepInjectionConfig      `json:"sleep,omitempty"`
	Mock       *MockInjectionConfig       `json:"mock"`
	Breakpoint *BreakpointInjectionConfig `json:"breakpoint"`
}

type Type string

const (
	TypeOff        Type = ""
	TypeSleep      Type = "sleep"
	TypeMock       Type = "mock"
	TypeBreakpoint Type = "breakpoint"
)
