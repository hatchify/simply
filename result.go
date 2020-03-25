package simply

// Status represents the current state of the test and result
type Status int

const (
	// PendingTarget status means test has not yet set a target value
	PendingTarget = iota
	// PendingComparison status means test has not yet compared the expected value
	PendingComparison

	// FailPendingValidation status means the test has failed but not yet validated
	FailPendingValidation
	// PassPendingValidation status means the test passed but not yet validated
	PassPendingValidation

	// Failed means test was validated and failed
	Failed
	// Passed means test was validated and succeeded
	Passed
)

// testResult represents a test's completion status
type testResult struct {
	status Status

	output string
}
