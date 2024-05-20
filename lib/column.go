package lib

/*
Functions for managing column data
*/

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

/* ************************************************************************** */

func Create(args []string) {
	if 1 < len(args) {
		fmt.Fprintf(os.Stderr, ERR_MSG_CREATE_ARGS)
	} else {
		if len(args[0]) == 0 { //no arg, create row
			CreateRow(app_data.Data)
		} else { //create column
			app_data.Data = CreateColumn(app_data.Data, args[0])
		}
	}
}

func CreateColumn(data DataBase, column string) DataBase {
	size := data_length()
	if data.Columns == nil {
		data.Columns = make(map[string][]interface{})
	}
	if size < 1 {
		size = 1
	}
	data.Columns[column] = make([]interface{}, 0, size)
	for i := 0; i < size; i++ {
		data.Columns[column] = append(data.Columns[column], 0.0)
	}
	return data
}

/* Create a row of zeros */
func CreateRow(data DataBase) {
	column_data := data.Columns
	for k, v := range column_data {
		data.Columns[k] = append(v, 0.0)
	}
}

/* ************************************************************************** */

func Read(args []string) {
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, ERR_MSG_READ_ARGS)
	} else {
		column := args[0]
		row, err := strconv.Atoi(args[1])
		if err == nil {
			ReadCell(app_data.Data, column, row)
		}
	}
}

// Read a specific value from the column table ; called with read command
func ReadCell(data DataBase, key string, row int) {
	if data.Columns[key] == nil {
		fmt.Fprintf(os.Stderr, ERR_MSG_COL_NOT_FOUND, key)
	} else {
		max := len(data.Columns[key]) - 1
		if max < row || row < 0 {
			fmt.Fprintf(os.Stderr, ERR_MSG_ROW_BETWEEN, max)
		} else {
			data := data.Columns[key][row]
			fmt.Printf("%s[%d]=%+v\n", key, row, data)
		}
	}
}

/* ************************************************************************** */

func Update(args []string) {
	if len(args) <= 2 {
		fmt.Fprintf(os.Stderr, ERR_MSG_UPDATE_ARGS)
	} else {
		column := args[0]
		row, row_err := strconv.Atoi(args[1])
		value := strings.Join(args[2:], " ")
		if row_err == nil {
			UpdateCell(app_data.Data, column, row, value)
		}
	}
}

// Update a specific value from the column table
func UpdateCell(data DataBase, key string, row int, value string) {
	column := data.Columns[key]
	if column == nil {
		fmt.Fprintf(os.Stderr, ERR_MSG_COL_NOT_FOUND, key)
	} else {
		max := len(column) - 1
		if row < 0 || max < row {
			fmt.Fprintf(os.Stderr, ERR_MSG_ROW_BETWEEN, max)
		} else {
			//if value can be turned into a number, then stuff it as a number
			if number, err := strconv.ParseFloat(value, 64); err == nil {
				//no error, value is a number
				data.Columns[key][row] = number
			} else {
				data.Columns[key][row] = value
			}
		}
	}
}

/* ************************************************************************** */

func Delete(args []string) {
	if len(args) != 1 || len(args[0]) < 1 {
		fmt.Fprintf(os.Stderr, ERR_MSG_DELETE_ARGS)
	} else {
		row, err := strconv.Atoi(args[0])
		if err == nil {
			DeleteRow(app_data.Data, row)
		} else { //delete column
			DeleteColumn(app_data.Data, args[0]) //TODO: add way to delete column
		}
	}
}

// Delete a row from all columns
func DeleteRow(data DataBase, row int) {
	for k, v := range data.Columns {
		max := len(v) - 1
		//while we have the first column, check the length before going on
		if max < row || row < 0 {
			fmt.Fprintf(os.Stderr, ERR_MSG_ROW_BETWEEN, max)
			break
		} else {
			copy(v[row:], v[row+1:])
			v[len(v)-1] = ""
			v = v[:len(v)-1]
			data.Columns[k] = v
		}
	}
}

func DeleteColumn(data DataBase, column string) {
	delete(data.Columns, column)
}
