package simply

import (
	"fmt"
	"path"
	"runtime"
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

	callingLine string
	callingFunc string

	result Result
}

var funcs = map[string]interface{}{}

// Test returns a new instance of a test
func Test(context *testing.T, name string) (test *Simply) {
	var t Simply
	t.Name = name
	t.context = context
	t.Validate = context.Error
	t.result.Status = PendingExpects

	if testing.Verbose() {
		if _, abs, line, ok := runtime.Caller(1); ok {
			_, file := path.Split(abs)
			t.callingLine = fmt.Sprintf("%s:%d", file, line)
			t.callingFunc = fmt.Sprintf("%s::%s", context.Name(), name)

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
func Target(target interface{}, context *testing.T, name string) (test *Simply) {
	t := Test(context, name)
	t.target = stringify(target)
	t.result.Status = PendingComparison

	return t
}

// GetStatus returns the current state of the test execution
func (t *Simply) GetStatus() Status {
	return t.result.Status
}

// Target begins a test with a target value to check
// Returns a reference to test for chaining convenience
func (t *Simply) Target(target interface{}) *Simply {
	t.result.Status = PendingComparison
	t.target = stringify(target)
	return t
}

// Assert requires the result, or stops the test
func (t *Simply) Assert() *Simply {
	t.Validate = t.context.Fatal
	return t
}

// Convenience for syntax sugar
func Assert(t *Simply) *Simply {
	return t.Assert()
}

// Equals ends a test with an expected value to compare to target
// Returns a reference to test Result for validation
func (t *Simply) Equals(val interface{}) *Result {
	t.expected = stringify(val)
	if t.target == t.expected {
		t.handlePass()

	} else {
		msg := fmt.Sprintf("%s - Failed! Expected <%+v> but got: <%+v>", t.Name, t.expected, t.target)
		t.handleFail(msg)
	}

	return &t.result
}

// DoesNotEqual ends a test with an expected value to compare to target
// Returns a reference to test Result for validation
func (t *Simply) DoesNotEqual(val interface{}) *Result {
	t.expected = stringify(val)
	if t.target != t.expected {
		t.handlePass()

	} else {
		msg := fmt.Sprintf("%s - Failed! Expected not to equal: <%+v>", t.Name, t.expected)
		t.handleFail(msg)
	}

	return &t.result
}

func (t *Simply) handlePass() {
	// Set validation handler to success
	// t.Validate() should print success output :)
	t.Validate = t.reportSuccess

	t.result.Status = PassPendingValidation
	t.result.output = fmt.Sprintf("    %s - Passed!", t.Name)
}

func (t *Simply) handleFail(msg string) {
	// Validation handler should already be set to t.context.Error()
	// t.Validate() should handle failures using default stdlib testing

	t.result.Status = FailPendingValidation
	t.result.output = msg
}

func (t *Simply) reportSuccess(a ...interface{}) {
	if testing.Verbose() {
		fmt.Println(t.result)
	}
}
