package injection

func dummyStatement(injectionName string, command BreakpointCommand) {} //  <---  set breakpoint here

func Breakpoint(injectionName string, config BreakpointInjectionConfig) {
	dummyStatement(injectionName, config.Command)
}

type BreakpointInjectionConfig struct {
	Command BreakpointCommand `json:"command" yaml:"command"`
}

type BreakpointCommand string

const (
	BreakpointManualInterruptType BreakpointCommand = "manual_interrupt"
)
