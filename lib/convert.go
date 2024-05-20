package lib

import (
	"fmt"
	"strconv"
)

//MARK - String

func is_interface_a_string(raw interface{}) bool {
	ret := false
	switch raw.(type) {
	case string:
		ret = true
	default:
		ret = false
	}
	return ret
}

func Interface_to_string(raw interface{}) string {
	return interface_to_string(raw)
}

func interface_to_string(raw interface{}) string {
	ret := ""
	switch i := raw.(type) {
	case string:
		ret = i
	case float64:
		ret = fmt.Sprintf("%f", i)
	case float32:
		ret = fmt.Sprintf("%f", i)
	case int64:
		ret = fmt.Sprintf("%0.0d", i)
	case int32:
		ret = fmt.Sprintf("%0.0d", i)
	case int:
		ret = fmt.Sprintf("%0.0d", i)
	default:
		fmt.Printf("got here")
	}
	return ret
}

func String_to_interface(raw string) interface{} {
	return string_to_interface(raw)
}

func string_to_interface(raw string) interface{} {
	var converted interface{}
	if number, err := strconv.ParseInt(raw, 10, 32); err == nil {
		fmt.Printf("here with int %d.\n", number)
		converted = int(number)
	} else if number, err := strconv.ParseFloat(raw, 64); err == nil {
		converted = number
	} else if state, err := strconv.ParseBool(raw); err == nil {
		converted = state
	} else {
		converted = raw
	}
	return converted
}

//MARK - Number

func is_interface_a_number(raw interface{}) bool {
	ret := false
	switch raw.(type) {
	case string:
		ret = false
	case float64:
		ret = true
	case float32:
		ret = true
	case int64:
		ret = true
	case int32:
		ret = true
	case int:
		ret = true
	default:
		ret = false
	}
	return ret
}

func Interface_to_float(raw interface{}) float64 {
	return interface_to_float(raw)
}

func interface_to_float(raw interface{}) float64 {
	ret := 0.0
	switch i := raw.(type) {
	case float64:
		ret = float64(i)
	case float32:
		ret = float64(i)
	case int64:
		ret = float64(i)
	}
	return ret
}

//MARK - Boolean

func is_interface_a_bool(raw interface{}) bool {
	ret := false
	switch raw.(type) {
	case bool:
		ret = true
	default:
		ret = false
	}
	return ret
}

func Interface_to_bool(raw interface{}) bool {
	return interface_to_bool(raw)
}

func interface_to_bool(raw interface{}) bool {
	ret := false
	switch i := raw.(type) {
	case bool:
		ret = i
	case string:
		ret = i == "true"
	case float64, float32:
		ret = i != 0.0
	case int64, int32, int:
		ret = i != 0
	default:
	}
	return ret
}
