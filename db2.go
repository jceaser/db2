package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jceaser/db2/lib"
	"github.com/peterh/liner"
)

var (
	history_fn = filepath.Join("db2_history")                   //used by liner
	names      = []string{"Create", "Read", "Update", "Delete"} //used by liner
)

/* Create a base AppData application with defined defaults */
func init_app() lib.AppData {
	var app_data = lib.AppData{
		Backlog_command: "",
		Worker_command:  "",
		IndentFile:      true,
		Verbose:         false,
		Format:          lib.CreateFormat(),
		Sort:            true}
	return app_data
}

/* Find the path to the history file for use with liner */
func history_path() string {
	history_base := ""
	if base, err := os.UserCacheDir(); err == nil {
		history_base = base
	} else if base, err = os.UserHomeDir(); err == nil {
		history_base = base
	} else {
		history_base = os.TempDir()
	}
	return history_base + "/" + history_fn
}

/* Setup the prompt reader */
func setup_liner(line *liner.State) {
	line.SetCtrlCAborts(true)

	line.SetTabCompletionStyle(liner.TabPrints)
	line.SetCompleter(func(line string) (c []string) {
		for _, n := range names {
			fmt.Print(n)
			if strings.HasPrefix(n, strings.ToLower(line)) {
				c = append(c, n)
				fmt.Print(n)
			}
		}
		return
	})
	if f, err := os.Open(history_path()); err == nil {
		line.ReadHistory(f)
		f.Close()
	}
}

/*
Run the interactive mode using the third party readline library. Help the
library store history, take each line and send it to ProcessLine()
*/
func InteractiveAdvance(line *liner.State, app_data *lib.AppData) lib.DataBase {
	fmt.Printf("Database by thomas.cherry@gmail.com\n")
	data := app_data.Data
	app_data.Running = true
	for app_data.Running == true {
		if name, err := line.Prompt(">"); err == nil {
			input := strings.Trim(name, " ")
			line.AppendHistory(name)
			lib.ProcessManyLines(input, app_data)
		} else if err == liner.ErrPromptAborted {
			fmt.Print("Aborted")
		} else {
			fmt.Print("Error reading line: ", err)
		}
		//save the history
		if f, err := os.Create(history_path()); err != nil {
			fmt.Print("Error creating history file: ", err)
		} else {
			line.WriteHistory(f)
			f.Close()
		}
	}
	lib.PrintCtrOnOut(lib.ESC_CURSOR_ON)
	return data
}

func main() {
	flag.Usage = func() {
		fmt.Printf("Forms Database\n")
		fmt.Printf("By thomas.cherry@gmail.com\n\n")
		raw_app_name := os.Args[0]
		index_of_slash := strings.LastIndex(raw_app_name, "/") + 1
		app_name := raw_app_name[index_of_slash:]
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", app_name)
		flag.PrintDefaults()
	}

	verbose := flag.Bool("verbose", false, "verbose")
	file_name := flag.String("file", "data.json", "data file")
	one_command := flag.String("command", "", "Run one command and exit")
	init_command := flag.String("init", "", "Run initial command and stay open")
	rpn_command := flag.String("rpn", "rpn", "command to process calculations")
	flag.Parse()

	app_data := init_app()

	app_data.Verbose = *verbose
	app_data.ActiveFile = *file_name
	app_data.Rpn = *rpn_command

	data := lib.Load(app_data.ActiveFile)

	if data == nil {
		fmt.Printf("Could not load data\n")
		os.Exit(1)
	} else {
		fmt.Printf("Data loaded from %s.\n", app_data.ActiveFile)
		app_data.Data = *data
	}
	lib.SetAppData(app_data)
	if 0 < len(*one_command) {
		lib.ProcessManyLines(*one_command, &app_data)
	} else {
		if 0 < len(*init_command) {
			lib.ProcessManyLines(*init_command, &app_data)
		}
		//readline setup
		line := liner.NewLiner()
		defer line.Close()
		setup_liner(line)

		//h := int(getHeight())
		//w := int(getWidth())

		if app_data.Verbose {
			lib.List(app_data.Data)
		}
		InteractiveAdvance(line, &app_data)
		if app_data.Verbose {
			lib.List(app_data.Data)
		}
	}
}
