package verifier

type Verifier interface {
	// Assume agent writes logs, we analyze them.
	// Error != nil signals problem for further fuzzing
	Verify(logsPath string) error
}
