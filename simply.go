package simply

import (
	"encoding/json"
	"fmt"
	"path"
	"runtime"
	"strconv"
	"testing"
)

type comparable string

// Simply is a test context returned from simply.Test(*testing.T, Name)
type Simply struct {
	name     string
	target   comparable
	expected comparable
	result   string

	AndValidate func(args ...interface{})

	callingLine string
	callingFunc string
}

// Test returns a new instance of a test and a pointer to the eventual result
func Test(context *testing.T, name string) (test *Simply, result *string) {
	var t Simply
	t.name = name
	t.AndValidate = context.Error

	if _, abs, line, ok := runtime.Caller(1); ok {
		_, file := path.Split(abs)
		t.callingLine = fmt.Sprintf("%s:%d", file, line)
		t.callingFunc = fmt.Sprintf("%s::%s", context.Name(), name)
	}

	return &t, &t.result
}

func (t *Simply) Expects(target interface{}) *Simply {
	t.target = stringify(target)
	return t
}

func stringify(from interface{}) (comp comparable) {
	if val, ok := from.(string); ok {
		return comparable(val)
	}

	if val, ok := from.(int); ok {
		return comparable(strconv.Itoa(val))
	}

	if s, err := json.MarshalIndent(from, "", " "); err == nil {
		return comparable(s)
	}

	panic("Unable to compare value")
}

func (t *Simply) success(a ...interface{}) {
	fmt.Println(t.result)
}

func (t *Simply) ToEqual(val interface{}) *Simply {
	t.expected = stringify(val)

	if t.target == t.expected {
		t.AndValidate = t.success
		t.result = fmt.Sprintf("    %s: %s - Passed!", t.callingLine, t.callingFunc)
	} else {
		t.result = fmt.Sprintf("%s - Failed! Expected <%+v> but got: <%+v>", t.callingFunc, t.expected, t.target)
	}

	return t
}
