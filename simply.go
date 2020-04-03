package simply

import (
	"fmt"
	"path"
	"runtime"
	"strings"
	"testing"
)

type comparable string

// Simply is a test context returned from simply.Test(Name, *testing.T)
type Simply struct {
	// Name provided at test initialization
	Name string

	// Test.Validate() parses the test Result and reports success/failure
	Validate func(args ...interface{})

	context *testing.T

	target   comparable
	expected comparable

	callingFile string
	callingFunc string
	callingLine string

	assert bool
	result *testResult
}

var funcs = map[string]interface{}{}

// Test returns a new instance of a test
func Test(context *testing.T, name string) (test *Simply) {
	var t Simply
	t.Name = "\"" + name + "\""
	t.context = context
	t.Validate = t.context.Error

	var result testResult
	t.result = &result
	t.result.status = PendingTarget

	if _, abs, line, ok := runtime.Caller(2); ok {
		_, file := path.Split(abs)
		t.callingFile = file
		t.callingLine = fmt.Sprintf("%s:%d", file, line)
		t.callingFunc = fmt.Sprintf("%s::%s", context.Name(), name)

		if testing.Verbose() {
			if _, ok := funcs[context.Name()]; !ok {
				// Space between tests
				if len(funcs) > 0 {
					fmt.Println("")
				}

				// Only print test name once
				funcs[context.Name()] = nil
				fmt.Println(context.Name())
			}
		}
	}

	return &t
}

// Target returns a new instance of a test with target set
func Target(target interface{}, context *testing.T, name string) *Simply {
	t := Test(context, name)
	t.target = makeComparable(target)
	t.result.status = PendingComparison

	return t
}

// GetStatus returns the current state of the test execution
func (t *Simply) GetStatus() Status {
	return t.result.status
}

// Target begins a test with a target value to check
// Returns a reference to test for chaining convenience
func (t *Simply) Target(target interface{}) *Simply {
	t.result.status = PendingComparison
	t.target = makeComparable(target)
	return t
}

// Assert requires the result, or stops the test
func (t *Simply) Assert() *Simply {
	t.assert = true
	return t
}

// Assert is static convenience for syntax sugar
func Assert(t *Simply) *Simply {
	return t.Assert()
}

// Equals ends a test with an expected value to compare to target
// Returns a reference to test Result for validation
func (t *Simply) Equals(val interface{}) *Simply {
	t.expected = makeComparable(val)
	if t.target == t.expected {
		msg := ""
		msg += fmt.Sprintf("Pass :: %s - [Equals]", t.Name)
		t.handlePass(msg)
	} else {
		msg := fmt.Sprintf("Fail :: %s - [Equals]", t.Name)
		if testing.Verbose() {
			msg = fmt.Sprintf("Fail :: %s - [Equals] - Expected <%+v> but got: <%+v>", t.Name, t.expected, t.target)
		}
		t.handleFail(msg)
	}

	return t
}

// DoesNotEqual ends a test with an expected value to compare to target
// Returns a reference to test Result for validation
func (t *Simply) DoesNotEqual(val interface{}) *Simply {
	t.expected = makeComparable(val)
	if t.target != t.expected {
		msg := ""
		msg += fmt.Sprintf("Pass :: %s - [DoesNotEqual]", t.Name)
		t.handlePass(msg)

	} else {
		msg := fmt.Sprintf("Fail :: %s - [DoesNotEqual]", t.Name)
		if testing.Verbose() {
			msg = fmt.Sprintf("Fail :: %s - [DoesNotEqual] - Expected not to equal: <%+v>", t.Name, t.expected)
		}
		t.handleFail(msg)
	}

	return t
}

func (t *Simply) handlePass(msg string) {
	// Set validation handler to success
	// t.Validate() should print success output :)

	if len(t.result.output) > 0 {
		t.result.output += "\n        "
	} else {
		t.result.output += "\n"
	}

	switch t.result.status {
	case PendingComparison, PassPendingValidation:
		t.result.output += msg

		t.result.status = PassPendingValidation
		t.Validate = t.reportSuccess

	case Passed, Failed:
		t.handlePostValidationError()
	}
}

func (t *Simply) handleFail(msg string) {
	// Validation handler should already be set to t.context.Error()
	// t.Validate() should handle failures using default stdlib testing

	t.result.output += "\n"

	switch t.result.status {
	case PendingComparison, PassPendingValidation, FailPendingValidation:
		t.result.output += msg

		t.result.status = FailPendingValidation
		if t.assert {
			t.Validate = t.context.Fatal
		} else {
			t.Validate = t.context.Error
		}

	case Passed, Failed:
		t.handlePostValidationError()
	}
}

func (t *Simply) handlePostValidationError() {
	t.result.status = Failed
	t.result.output = fmt.Sprintf("\nError :: %s - [Test Sequence] - ", t.Name)
	if _, _, line, ok := runtime.Caller(2); ok {
		t.result.output += fmt.Sprintf("%s:%d:\n  -- Please avoid running comparisons after validation", t.callingFile, line)
	} else {
		t.result.output += t.callingLine + "\n  -- Please avoid running comparisons after validation"
	}
	t.context.Error(t.result.output)
	t.Validate = t.duplicateValidation

}

func (t *Simply) reportSuccess(a ...interface{}) {
	t.result.status = Passed
	t.Validate = t.duplicateValidation

	if testing.Verbose() {
		fmt.Println("        " + strings.TrimSpace(t.String()))
	}
}

func (t *Simply) duplicateValidation(a ...interface{}) {
	t.result.status = Failed

	t.result.output = fmt.Sprintf("\nError :: %s - [Test Sequence] - ", t.Name)
	if _, _, line, ok := runtime.Caller(1); ok {
		t.result.output += fmt.Sprintf("%s:%d:\n  -- Please avoid running validation on the same test twice", t.callingFile, line)
	} else {
		t.result.output += t.callingLine + "\n  -- Please avoid running validation on the same test twice"
	}
	t.context.Error(t.result.output)
}

func (t *Simply) String() string {
	switch t.result.status {
	case PassPendingValidation, Passed:
		// Test just validated and passed
		return t.result.output
	case FailPendingValidation, Failed:
		// Test just validated and failed
		t.result.status = Failed
		t.Validate = t.duplicateValidation
		return t.result.output
	case PendingTarget:
		return "\nError :: " + t.Name + " - [Test Sequence]\n  -- Please begin test with simply.Target() or add target to test with test.Target() before comparison or validation"
	case PendingComparison:
		return "\nError :: " + t.Name + " - [Test Sequence]\n  -- Please perform comparison with test.Equals() or target.Equals() before validation"
	}

	return t.result.output
}
