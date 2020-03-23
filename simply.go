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

	target   comparable
	expected comparable

	callingLine string
	callingFunc string

	result Result
}

var funcs = map[string]interface{}{}

// Test returns a new instance of a test
func Test(name string, context *testing.T) (test *Simply) {
	var t Simply
	t.Name = name
	t.Validate = context.Error
	t.result.Status = PendingExpects

	if _, abs, line, ok := runtime.Caller(1); ok {
		_, file := path.Split(abs)
		t.callingLine = fmt.Sprintf("%s:%d", file, line)
		t.callingFunc = fmt.Sprintf("%s::%s", context.Name(), name)

		if _, ok := funcs[context.Name()]; !ok {
			funcs[context.Name()] = nil
			fmt.Println("\n" + context.Name())
		}
	}

	return &t
}

// TestTarget returns a new instance of a test with target set
func TestTarget(target interface{}, name string, context *testing.T) (test *Simply) {
	t := Test(name, context)
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

// Equals ends a test with an expected value to compare to target
// Returns a reference to test Result for validation
func (t *Simply) Equals(val interface{}) *Result {
	t.expected = stringify(val)
	if t.target == t.expected {
		t.handlePass()

	} else {
		t.handleFail()
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

func (t *Simply) handleFail() {
	// Validation handler should already be set to t.context.Error()
	// t.Validate() should handle failures using default stdlib testing

	t.result.Status = FailPendingValidation
	t.result.output = fmt.Sprintf("%s - Failed! Expected <%+v> but got: <%+v>", t.Name, t.expected, t.target)
}

func (t *Simply) reportSuccess(a ...interface{}) {
	fmt.Println(t.result)
}
