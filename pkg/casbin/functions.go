package casbin

import (
	"fmt"
	"strconv"
	"time"
)

func SinceFunc(args ...any) (any, error) {
	if len(args) != 2 {
		return false, fmt.Errorf("expected 2 values, got %d", len(args))
	}

	last, isTime := args[0].(time.Time)
	if !isTime {
		return false, fmt.Errorf("expected time.Time, got %T", args[0])
	}

	delay, err := getFloat64Value(args[1])
	if err != nil {
		return false, fmt.Errorf("delay is not a numeric")
	}

	return time.Since(last) > time.Duration(delay), nil
}

func InRangeFunc(args ...any) (any, error) {
	if len(args) != 4 {
		return false, fmt.Errorf("expected 4 values, got %d", len(args))
	}

	rBottomValue, err := getFloat64Value(args[0])
	if err != nil {
		return false, fmt.Errorf("request bottom value: %w", err)
	}

	pBottomValue, err := getFloat64Value(args[1])
	if err != nil {
		return false, fmt.Errorf("policy bottom value: %w", err)
	}

	rTopValue, err := getFloat64Value(args[2])
	if err != nil {
		return false, fmt.Errorf("request top value is not a float64")
	}

	pTopValue, err := getFloat64Value(args[3])
	if err != nil {
		return false, fmt.Errorf("policy top value: %w", err)
	}

	return rBottomValue >= pBottomValue && rTopValue <= pTopValue, nil
}

func getFloat64Value(value any) (float64, error) { //nolint:gocyclo
	switch tValue := value.(type) {
	case int:
		return float64(tValue), nil
	case int8:
		return float64(tValue), nil
	case int16:
		return float64(tValue), nil
	case int32:
		return float64(tValue), nil
	case int64:
		return float64(tValue), nil
	case uint:
		return float64(tValue), nil
	case uint8:
		return float64(tValue), nil
	case uint16:
		return float64(tValue), nil
	case uint32:
		return float64(tValue), nil
	case uint64:
		return float64(tValue), nil
	case float32:
		return float64(tValue), nil
	case float64:
		return tValue, nil
	case string:
		floatValue, err := strconv.ParseFloat(tValue, 64)
		if err != nil {
			return 0, fmt.Errorf("parse float64: %w", err)
		}

		return floatValue, nil
	default:
		return 0, fmt.Errorf("unsupported type %T", tValue)
	}
}
