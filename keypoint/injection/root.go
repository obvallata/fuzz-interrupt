package injection

type Config struct {
	Type Type `json:"type" yaml:"type"`

	Sleep *SleepInjectionConfig `json:"sleep,omitempty" yaml:"sleep"`
	Mock  *MockInjectionConfig  `json:"mock,omitempty" yaml:"mock"`
}

type Type string

const (
	TypeOff        Type = ""
	TypeSleep      Type = "sleep"
	TypeMock       Type = "mock"
	TypeBreakpoint Type = "breakpoint"
)

func NewDefaultOffConfig() Config {
	return Config{Type: TypeOff}
}
