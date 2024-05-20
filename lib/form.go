package lib

/*
Functions for managing forms
*/

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/peterh/liner"
)

func FormCreate(args []string) {
	if len(args) < 2 {
		e(ERR_MSG_FORM_create)
	} else {
		name := arg(args, 0, "")
		if len(name) < 1 {
			e(ERR_MSG_FORM_REQUIRED, name)
		} else {
			if app_data.Data.Forms[name] != nil {
				e(ERR_MSG_FORM_EXISTS, name)
			} else {
				items := args[1:]
				app_data.Data.Forms[name] = items
			}
		}
	}
}

func FormRead(args []string) {
	name := arg(args, 0, "")
	if len(name) < 1 {
		fmt.Printf("%+v\n", app_data.Data.Forms) //TODO: make this pretty
	} else {
		fmt.Printf("%+v\n", app_data.Data.Forms[name]) //TODO: make this pretty
	}
}

func FormUpdate(args []string) {
	if len(args) < 2 {
		e(ERR_MSG_FORM_UPDATE)
	} else {
		name := arg(args, 0, "")
		if len(name) < 1 {
			e(ERR_MSG_FORM_REQUIRED)
		} else {
			items := args[1:]
			if app_data.Data.Forms[name] == nil {
				e(ERR_MSG_FORM_NOT_EXISTS, name)
			} else {
				app_data.Data.Forms[name] = items
			}
		}
	}
}

func FormDelete(args []string) {
	name := arg(args, 0, "")
	if len(name) < 1 {
		e(ERR_MSG_FORM_REQUIRED)
		return
	}
	delete(app_data.Data.Forms, name)
}

func FormRename(args []string) {
	if len(args) < 2 {
		e(ERR_MSG_FORM_RENAME)
	} else {
		src_name := arg(args, 0, "")
		dest_name := arg(args, 1, "")
		if 0 < len(src_name) && 0 < len(dest_name) {
			app_data.Data.Forms[dest_name] = app_data.Data.Forms[src_name]
			delete(app_data.Data.Forms, src_name)
		}
	}
}

//MARK -

/** A simple form filler to create a new row using a form as input */
func _FormFiller(form string, action string) {
	dry_run := false
	if action == "dry-run" {
		dry_run = true
	}
	line := liner.NewLiner()
	defer line.Close()
	if !dry_run {
		CreateRow(app_data.Data) //new row
	}
	row := data_length() - 1
	for _, column := range app_data.Data.Forms[form] {
		if 0 < len(app_data.Data.Calculations[column]) {
			continue // this is a calculation, skip it
		}
		var answer interface{}
		answer = 0.0
		fmt.Printf("Enter in a value for column '%s'.\n", column)
		raw_response, _ := line.Prompt("#")
		quiters := []string{"stop", "exit", "quit"}
		if contains(quiters, raw_response) {
			return
		}
		number, err := strconv.ParseFloat(raw_response, 64)
		if err == nil {
			answer = number
		} else {
			answer = raw_response
		}
		if !dry_run {
			app_data.Data.Columns[column][row] = answer
		}
	}
}

/* Use control codes to draw a form that the user then fills out */
func FormFillerVisual(form string, action string) {
	ScrSave()

	//Lines
	LINE_QUESTION := 1
	LINE_INPUT := 2
	LINE_MESSAGE := 3
	LINE_FORM_START := 5
	LINE_FORM_HEAD := LINE_FORM_START - 1

	temp_data := make(map[string]interface{})
	reviewing := true // main loop flag
	dry_run := false  // dry run mode, ask for input, but do not save
	mode := "normal"

	/** some inner utility functions for this function */

	/* * * * * * * * * * * * * * * * */
	/** run calculations from the database and save to a temp map */
	populate_calcs := func() {
		var json_data = DataBase{}
		for k, v := range temp_data {
			json_data = CreateColumn(json_data, k)
			json_data.Columns[k][0] = v
		}
		for key, formula := range app_data.Data.Calculations {
			result := formula_for_data(formula, 0, json_data)
			result_as_float, _ := strconv.ParseFloat(result, 64)
			json_data = CreateColumn(json_data, key)
			json_data.Columns[key][0] = result_as_float
			temp_data[key] = result_as_float
		}
	}

	/* * * * * * * * * * * * * * * * */
	/** copy in values from the database to a temp map and also calculations */
	populate_temp := func(index int) {
		for _, column := range app_data.Data.Forms[form] {
			if 0 < len(app_data.Data.Calculations[column]) {
				continue // this is a calculation, skip it
			}
			temp_data[column] = app_data.Data.Columns[column][index]
		}
		populate_calcs()
	}

	/* * * * * * * * * * * * * * * * */
	/** ask the user for an answer to a question, allow for quit commands */
	asker := func(line *liner.State, question, suggestion string, options ...interface{}) (string, string) {
		PrintStrAt(fmt.Sprintf(question, options...), LINE_QUESTION, 1)

		PrintCtrAt(ESC_CLEAR_LINE, LINE_INPUT, 1)
		PrintStrAt(fmt.Sprintf(""), LINE_INPUT, 1)
		raw_response, _ := line.PromptWithSuggestion("#", suggestion, len(suggestion))
		PrintCtrAt(ESC_CLEAR_LINE, LINE_INPUT, 1)

		var none string

		quiters := []string{"stop", "exit", "q", "quit"}
		if contains(quiters, raw_response) {
			ScrRestore()
			return none, "exit"
		}
		return raw_response, none

	}

	/* * * * * * * * * * * * * * * * */
	/** draw the values in the form */
	draw_values := func() {
		//setup
		PrintStrAt("", LINE_FORM_HEAD, 1)
		fmt.Print(strings.Repeat(fmt.Sprintf("%c", RuneS3), 80))
		for c_count, column := range app_data.Data.Forms[form] {
			answer := temp_data[column]
			if answer == nil {
				answer = "<empty>"
			}
			PrintCtrAt(ESC_CLEAR_LINE, c_count+LINE_FORM_START, 1)
			field := fmt.Sprintf("%d: %s = %v.\n", c_count, column, answer)
			PrintStrAt(field, c_count+LINE_FORM_START, 1)
			Table("")
		}
	}

	/**********************************/

	line := liner.NewLiner()

	wc := func(line string, pos int) (head string, completions []string, tail string) {
		fmt.Print("\a")
		return "", []string{""}, ""
	}
	line.SetWordCompleter(wc)

	defer line.Close()

	if action == "dry-run" {
		dry_run = true
	} else if action == "show" {
		//iterate over all the values
		mode = action
		max := DataLength(app_data.Data)
		index := 0
		for reviewing {
			populate_temp(index)
			draw_values()
			result, exit_command := asker(line,
				"%s %d of %d. Type pre, next, or quit.",
				"next",
				mode,
				index+1,
				max)
			if exit_command != "" {
				break
			} else if strings.HasPrefix("next", result) {
				index = int(math.Min(float64(index+1), float64(max-1)))
			} else if strings.HasPrefix("previous", result) {
				index = int(math.Max(0.0, float64(index-1)))
			} else if strings.HasPrefix("delete", result) {
				DeleteRow(app_data.Data, index)
			}
		}
		ScrRestore()
		return
	}

	for reviewing {
		draw_values()

		//review and for loop were here
		for c_count, column := range app_data.Data.Forms[form] {
			if 0 < len(app_data.Data.Calculations[column]) {
				continue // this is a calculation, skip it
			}
			asking := true
			var answer interface{}
			answer = 0.0
			for asking {
				raw_response, quiter := asker(line,
					"Enter in a number for column '%s'.\n",
					Green(column))
				if quiter != "" {
					ScrRestore()
					return
				}

				number, err := strconv.ParseFloat(raw_response, 64)
				if err != nil {
					PrintCtrAt(ESC_CLEAR_LINE, LINE_MESSAGE, 1)
					answer = raw_response
					asking = false
				} else {
					PrintCtrAt(ESC_CLEAR_LINE, LINE_MESSAGE, 1)
					answer = number
					asking = false
				}
			}
			PrintCtrAt(ESC_CLEAR_LINE, c_count+LINE_FORM_START, 1)
			msg := fmt.Sprintf("%d: %s = %f\n", c_count, Green(column), answer)
			PrintStrAt(msg, c_count+LINE_FORM_START, 1)
			var a interface{}
			a = answer
			temp_data[column] = a
		}

		populate_calcs()
		draw_values()

		PrintCtrAt(ESC_CLEAR_LINE, LINE_QUESTION, 1)
		raw_response, _ := line.Prompt("done? yes or no: ")
		PrintCtrAt(ESC_CLEAR_LINE, LINE_QUESTION, 1)

		quiters := []string{"stop", "exit", "e", "done", "d", "yes", "y", "save"}
		if contains(quiters, raw_response) {
			reviewing = false
			if !dry_run {
				CreateRow(app_data.Data) //new row
				row := data_length() - 1
				for k, v := range temp_data {
					_, okay := app_data.Data.Columns[k]
					if okay {
						//should also create the new row here
						app_data.Data.Columns[k][row] = v
					}
				}
			}
		}
	}
	ScrRestore()
}
