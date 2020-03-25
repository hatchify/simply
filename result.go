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

// Result represents a test's completion status
type Result struct {
	Status Status

	output string
}

func (r *Result) String() string {
	switch r.Status {
	case PassPendingValidation, Passed:
		// Test just validated and passed
		r.Status = Passed
		return r.output
	case FailPendingValidation, Failed:
		// Test just validated and failed
		r.Status = Failed
		return r.output
	}

	// If we don't have output
	return toString(*r)
}
