package lib

/*
AppData Structure and any related functions
*/

type AppData struct {
    Backlog_command string
    Worker_command string
    Backlog_list []string
    Rpn string
    ColumnCache map[string][]interface{}

    Data DataBase
    Verbose bool
    IndentFile bool
    ActiveFile string
    Running bool
    Format Format
    Sort bool
}
