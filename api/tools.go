package api

import (
	"strings"
	"time"
)

func EmptyValue(value interface{}) bool {
	if value == nil {
		return true
	}

	switch val := value.(type) {
	case []int32:
		return len(val) == 0
	case []int64:
		return len(val) == 0
	case []float32:
		return len(val) == 0
	case []float64:
		return len(val) == 0
	case []string:
		return len(val) == 0
	case int32:
		return val == 0
	case int64:
		return val == 0
	case float32:
		return val == 0
	case float64:
		return val == 0
	case time.Time:
		return val.IsZero()
	case *time.Time:
		if val == nil {
			return true
		} else {
			return val.IsZero()
		}
	case string:
		return strings.TrimSpace(val) == ""
	default:
		return false
	}
}
