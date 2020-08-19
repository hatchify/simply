package simply

import (
	"testing"
)

// Tests that the prehook runs before the test runs when target value is a function
func Test_FuncTarget(t *testing.T) {
	var preHookRan bool

	// Return true if the PreHook has run, false if not
	testFunc := func() bool {
		return preHookRan
	}

	testcases := []TestCase{
		{
			Message: "ensure prehook runs before target func",
			PreHook: func(tc TestCase) {
				preHookRan = true
			},
			Target:   testFunc(),
			Expected: true,
		},
	}

	Run(t, testcases)
}
