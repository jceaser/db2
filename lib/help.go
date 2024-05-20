package lib

import (
	"fmt"
	"strings"
)

/* print out a help method */
func Help() {
	fmt.Printf("Database by thomas.cherry@gmail.com\n")
	fmt.Printf("Manage table data with optional form display.\n")
	fmt.Printf("\nNote: Arguments with ? are optional\n\n")

	format := "%4s %-14s %-14s %-40s\n"

	forty := strings.Repeat("-", 40)
	fmt.Printf(format, "Flag", "Long", "Arguments", "Description")
	fmt.Printf(format, "----", "------------", "------------", forty)
	fmt.Printf(format, "c", "create", "name", "Create a column by name")
	fmt.Printf(format, "", "", "", "create a zero row")
	fmt.Printf(format, "r", "read", "col row", "read a column row")
	fmt.Printf(format, "u", "update", "col row val", "update a column row")
	fmt.Printf(format, "d", "delete", "index", "delete a row by number")
	fmt.Printf(format, "", "", "name", "delete a column by name")

	fmt.Printf(format, "a", "append", "<value list>", "append a table")
	fmt.Printf(format, "A", "append-by-name", "name:value...", "append a table with named columns")

	fmt.Printf(format, "n", "rename", "src dest", "rename a column from src to dest")
	fmt.Printf("\n")

	fmt.Printf(format, "fc", "form-create", "name list", "Create a form")
	fmt.Printf(format, "fr", "form-read", "name?", "list forms, all if name is not given")
	fmt.Printf(format, "fu", "form-update", "name formula", "update a form")
	fmt.Printf(format, "fd", "form-delete", "name", "delete a form")
	fmt.Printf(format, "fn", "form-rename", "src dest", "rename a form from src to dest")
	fmt.Printf(format, "ff", "form-filler", "name action?", "Create a row using a form")
	fmt.Printf(format, "FF", "Form-filler", "name action?", "Basic form filler")
	fmt.Printf("\n")

	fmt.Printf(format, "cc", "calc-create", "name formula", "Create a calculation")
	fmt.Printf(format, "cr", "calc-read", "name?", "list calculations, all if name is not given")
	fmt.Printf(format, "cu", "calc-update", "name formula", "update a calculation")
	fmt.Printf(format, "cd", "calc-delete", "name", "delete a calculation")
	fmt.Printf(format, "cn", "calc-rename", "src dest", "rename a calculation from src to dest")
	fmt.Printf("\n")

	fmt.Printf(format, "sc", "sum-create", "form list", "assign a list of sum functions")
	fmt.Printf(format, "sr", "sum-read", "form?", "read all values, or one sum function")
	fmt.Printf(format, "su", "sum-update", "form list", "update a list of sum functions")
	fmt.Printf(format, "sd", "sum-delete", "form", "clear a list of sum functions")
	fmt.Printf("\n")

	fmt.Printf(format, "sum", "summary", "form list", "summarize a form with function list:")
	fmt.Printf(format, "", "", "", "avg,count,max,min,medium,mode,min,nop,sum,sdev")
	fmt.Printf("\n")

	fmt.Printf(format, "t", "table", "form?", "display a table, optionally as a form")
	fmt.Printf(format, "l", "ls list", "", "Print columns with their data.")
	fmt.Printf(format, "", "row", "form row?", "show a header and one row from form")
	fmt.Printf("\n")

	fmt.Printf(format, "s", "save", "", "save database to file")
	fmt.Printf(format, "", "dump", "", "output the current data")
	fmt.Printf(format, "q", "quit", "", "quit interactive mode")
	fmt.Printf(format, "", "exit", "", "quit interactive mode")
	fmt.Printf(format, "h", "help", "", "this output")
	fmt.Printf(format, "e", "echo", "string", "echo out something")
	fmt.Printf(format, "-", "----", "sep count", "print out a separator")
	fmt.Printf(format, "", "file", "name?", "set or print current file name")
	fmt.Printf(format, "", "rpn-set", "path?", "set or print current rpn command")
	fmt.Printf(format, "", "verbose", "", "toggle verbose mode")
	fmt.Printf(format, "", "sort?", "", "output current sorting state")
	fmt.Printf(format, "", "sort", "", "toggle the current sort mode")
}
