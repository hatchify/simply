package simply

import (
	"encoding/json"
	"strconv"
)

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
