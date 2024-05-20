package main

/*
tests for the file convert.go
*/

import (
	"testing"

	"github.com/jceaser/db2/lib"
)

func TestConvert(t *testing.T) {
	//var data interface{}

	var f64 float64
	f64 = 3.145920
	ans := lib.Interface_to_string(f64)
	sline(t, "3.145920", ans, "Float64 to string does not equal %s from %s.")

	var f32 float32
	f32 = 3.145920
	ans = lib.Interface_to_string(f32)
	sline(t, "3.145920", ans, "Float32 to string does not equal %s from %s.")

	ans = lib.Interface_to_string(3.145920)
	sline(t, "3.145920", ans, "Float to string does not equal %s from %s.")

	ans = lib.Interface_to_string(1)
	sline(t, "1", ans, "Int to string does not equal %s from '%s'.")

	ans = lib.Interface_to_string("3.145920")
	sline(t, "3.145920", ans, "String to string does not equal %s from %s.")
}

func TestFloatConvert(t *testing.T) {
	ans := lib.Interface_to_float(3.14592)
	pline(t, 3.145920, ans, "Raw to float does not equal %f from %f.")

	ans = lib.Interface_to_float("3.14592")
	pline(t, 0.0, ans, "String test %f from %f.")

}

func TestBoolConvert(t *testing.T) {
	//bool tests
	are := func(expected bool, input interface{}, msg string) {
		ans := lib.Interface_to_bool(input)
		if ans != expected {
			t.Errorf(msg+"\n", expected, ans)
		}
	}

	//boolean tests
	are(true, true, "Raw to bool does not equal %t from %t.")
	are(false, false, "Raw to bool does not equal %t from %t.")
	//string tests
	are(true, "true", "String test %t from %t.")
	are(false, "false", "String test %t from %t.")
	are(false, "3.14", "String test %t from %t.")
	//float tests
	are(true, -3.14, "-Pi test %t from %t.")
	are(false, 0.0, "Zero float test %t from %t.")
	are(true, 3.14, "Pi test %t from %t.")
	//int tests
	are(true, -42, "-42 test %t from %t.")
	are(false, 0, "Zero test %t from %t.")
	are(true, 42, "42 test %t from %t.")
}

func TestStringBackToInterface(t *testing.T) {
	are := func(expected interface{}, input string, msg string) {
		actual := lib.String_to_interface(input)
		if actual != expected {
			t.Errorf("%s: %v != %v", msg, expected, actual)
		}
	}
	are("zero", "zero", "zero as a word")
	are(3.14, "3.14", "a float")
	are(3.14e+12, "3.14e+12", "a large float")
	are(8.589934592e+09, "8.589934592e+09", "a large int")
	are(42, "42", "an int")
	are(true, "true", "a true boolean")
	are(false, "false", "a false boolean")
}
