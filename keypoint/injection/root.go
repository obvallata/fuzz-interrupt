package injection

type Config struct {
	Type Type `json:"type" yaml:"type"`

	Sleep      *SleepInjectionConfig      `json:"sleep,omitempty" yaml:"sleep"`
	Mock       *MockInjectionConfig       `json:"mock,omitempty" yaml:"mock"`
	Breakpoint *BreakpointInjectionConfig `json:"breakpoint,omitempty" yaml:"breakpoint"`
}

type Type string

const (
	TypeOff        Type = ""
	TypeSleep      Type = "sleep"
	TypeMock       Type = "mock"
	TypeBreakpoint Type = "breakpoint"
)
