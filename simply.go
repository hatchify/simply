package simply

import (
	"fmt"
	"path"
	"runtime"
	"testing"
)

type comparable string

// Simply is a test context returned from simply.Test(*testing.T, Name)
type Simply struct {
	Name     string
	Validate func(args ...interface{})

	target   comparable
	expected comparable

	callingLine string
	callingFunc string

	result *Result
}

// Test returns a new instance of a test and a pointer to the eventual result
func Test(name string, context *testing.T) (test *Simply) {
	var t Simply
	t.Name = name
	t.Validate = context.Error

	if _, abs, line, ok := runtime.Caller(1); ok {
		_, file := path.Split(abs)
		t.callingLine = fmt.Sprintf("%s:%d", file, line)
		t.callingFunc = fmt.Sprintf("%s::%s", context.Name(), name)
	}

	return &t
}

func (t *Simply) Expects(target interface{}) *Simply {
	t.target = stringify(target)
	return t
}

func (t *Simply) success(a ...interface{}) {
	fmt.Println(t.result)
}

func (t *Simply) ToEqual(val interface{}) *Result {
	var result Result
	t.result = &result

	t.expected = stringify(val)
	if t.target == t.expected {
		// t.Validate should print success output :)
		t.result.output = fmt.Sprintf("    %s: %s - Passed!", t.callingLine, t.callingFunc)
		t.result.Success = true
		t.Validate = t.success
	} else {
		// t.Validate should equal *testing.T.Errof(), which handles failures using default stdlib testing
		t.result.output = fmt.Sprintf("%s - Failed! Expected <%+v> but got: <%+v>", t.callingFunc, t.expected, t.target)
	}

	t.result.Complete = true

	return t.result
}
