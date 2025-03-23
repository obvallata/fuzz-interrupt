package injection

func dummyStatement(keypointName string, command BreakpointCommand) {} //  <---  set breakpoint here

func Breakpoint(keypointName string, config BreakpointInjectionConfig) {
	dummyStatement(keypointName, config.Command)
}

type BreakpointInjectionConfig struct {
	Command BreakpointCommand `json:"command"`
}

type BreakpointCommand string

const (
	BreakpointNotifyStartType   BreakpointCommand = "notify_start"
	BreakpointNotifySuccessType BreakpointCommand = "notify_success"
	BreakpointNotifyErrorType   BreakpointCommand = "notify_error"
)
