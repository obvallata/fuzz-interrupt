package injection

func dummyStatement(injectionName string) {} //  <---  set breakpoint here

func Breakpoint(injectionName string) {
	dummyStatement(injectionName)
}
