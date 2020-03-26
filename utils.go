package simply

import (
	"encoding/json"
	"strconv"
)

func makeComparable(from interface{}) (comp comparable) {
	return comparable(Stringify(from))
}

func Stringify(from interface{}) string {
	if val, ok := from.(error); ok {
		return val.Error()
	}

	if val, ok := from.(int); ok {
		return strconv.Itoa(val)
	}

	if val, ok := from.(string); ok {
		return val
	}

	if s, err := json.MarshalIndent(from, "", " "); err == nil {
		return string(s)
	}

	panic("Unable to compare value")
}
