package main

/*
tests for the all files in the lib module
*/

import (
	"testing"

	"github.com/jceaser/db2/lib"
)

/******************************************************************************/

func init() {
	//no setup needed
}

/**************************************/

func TestAvergae(t *testing.T) {
	data := []interface{}{1.0, 2.0, 3.0, 10.0}
	ans := lib.Average(data)
	pline(t, 4.0, ans, "Average does not equal %f, got %f")
}
func TestMedian(t *testing.T) {
	data := []interface{}{1.0, 2.0, 3.0, 4.0, 5.0, 6.0}
	ans := lib.Median(data)
	pline(t, 3.5, ans, "Median does not %f, got %f")

	data = []interface{}{1.0, 2.0, 3.0}
	ans = lib.Median(data)
	pline(t, 2.0, ans, "Median does not equal %f, got %f")
}

func TestMode(t *testing.T) {
	data := []interface{}{0.0, 1.0, 2.0, 3.0, 3.0, 10.0}
	ans := lib.Mode(data)
	pline(t, 3.0, ans, "Mode does not equal %f")
}

func TestInitDataBase(t *testing.T) {
	data := lib.InitDataBase()
	msg := "init here"

	expected := []string{"foo", "bar", "foobar", "row"}
	actual := data.Forms["main"]

	for idx := 0; idx < len(actual); idx++ {
		if expected[idx] != actual[idx] {
			t.Errorf("%s, %v==%v\n", msg, expected, actual)
		}
	}
}

/*
table;append 1 2 3 ; table ; append 1.0 ; table ; append 5 4 3 2 1 ; table
*/
func TestAppend(t *testing.T) {
	data := lib.InitDataBase()
	//SetData(data)

	source := []string{"9", "8", "7"}
	expected := []string{"9.000000", "8.000000", "7.000000"}
	lib.AppendTable(data, source)
	row := lib.DataLength(data) - 1
	b := lib.Interface_to_string(data.Columns["bar"][row])
	f := lib.Interface_to_string(data.Columns["foo"][row])
	r := lib.Interface_to_string(data.Columns["rab"][row])
	ans := []string{b, f, r}
	check_three(t, expected, ans, "append test - exact - %s != expected[%d]=%s")

	source = []string{"8", "7", "6", "5"}
	expected = []string{"8.000000", "7.000000", "6.000000"}
	lib.AppendTable(data, source)
	row = lib.DataLength(data) - 1
	b = lib.Interface_to_string(data.Columns["bar"][row])
	f = lib.Interface_to_string(data.Columns["foo"][row])
	r = lib.Interface_to_string(data.Columns["rab"][row])
	ans = []string{b, f, r}
	check_three(t, expected, ans,
		"append test - to many given - %s != expected[%d]=%s")

	source = []string{"4"}
	expected = []string{"4.000000", "0.000000", "0.000000"}
	lib.AppendTable(data, source)
	row = lib.DataLength(data) - 1
	b = lib.Interface_to_string(data.Columns["bar"][row])
	f = lib.Interface_to_string(data.Columns["foo"][row])
	r = lib.Interface_to_string(data.Columns["rab"][row])
	ans = []string{b, f, r}
	check_three(t, expected, ans,
		"append test - not enough given - %s != expected[%d]=%s")

	source = []string{"foo:3.14"}
	expected = []string{"0.000000", "3.14", "0.000000"}
	lib.AppendTableByName(data, source)
	row = lib.DataLength(data) - 1
	b = lib.Interface_to_string(data.Columns["bar"][row])
	f = lib.Interface_to_string(data.Columns["foo"][row])
	r = lib.Interface_to_string(data.Columns["rab"][row])
	ans = []string{b, f, r}

	check_three(t, expected, ans,
		"append by name test - %s != expected[%d]=%s")
}

func TestCommands(t *testing.T) {
	data := lib.InitDataBase()
	lib.SetData(data)
	app_data := lib.GetAppData()

	//lib.ProcessManyLines("c name 0 ; u name \"test\"", app_data.data)
	lib.ProcessManyLines("create name", &app_data)

	expected1 := []string{"0.000000", "0.000000", "0.000000"}
	b1 := lib.Interface_to_string(data.Columns["name"][0])
	f1 := lib.Interface_to_string(data.Columns["name"][1])
	r1 := lib.Interface_to_string(data.Columns["name"][2])
	ans1 := []string{b1, f1, r1}
	check_three(t, expected1, ans1, "cmd test - create - %s != expected[%d]=%s")

	lib.ProcessManyLines("update name 1 10 ; update name 2 test", &app_data)

	expected := []string{"0.000000", "10.000000", "test"}
	b := lib.Interface_to_string(data.Columns["name"][0])
	f := lib.Interface_to_string(data.Columns["name"][1])
	r := lib.Interface_to_string(data.Columns["name"][2])
	ans := []string{b, f, r}
	check_three(t, expected, ans, "cmd test - update - %s != expected[%d]=%s")

}

/**************************************/

func length(data map[string][]interface{}) int {
	length := -1
	for _, v := range data {
		length = len(v)
		break
	}
	return length
}

func check_three(t *testing.T, expected []string, ans []string, msg string) {
	for i, v := range ans {
		if v != expected[i] {
			t.Errorf(msg, v, i, expected[i])
			break
		}
	}
}

/* String compare test */
func sline(t *testing.T, expected string, ans string, msg string) {
	if ans != expected {
		t.Errorf(msg+"\n", expected, ans)
	}
}

/* Float compare test */
func pline(t *testing.T, expected float64, ans float64, msg string) {
	if ans != expected {
		t.Errorf(msg+"\n", expected, ans)
	}
}
