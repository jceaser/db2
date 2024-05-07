package lib

/*
General Application code
*/

import (
	"fmt"
    "os"
    "bytes"
    "sort"
    "io/ioutil"
    "math"
    "strconv"
    "os/exec"
    "strings"
    "encoding/json"
)

type screen_buffers struct {
    left_hud string
    right_hud string
    content string
}

type Form struct {
    Name string
    Columns []string
    Settings map[string]string
}

var buffers = screen_buffers{left_hud: "", right_hud: "", content: ""}

var app_data AppData

const (
    ERR_MSG_COL_NOT_FOUND = "Column %s not found\n"
    ERR_MSG_ROW_BETWEEN = "Row must be between 0 and %d\n"
    ERR_MSG_VALUE_NUM = "Value '%s' is not a number\n"
    ERR_MSG_CREATE_ARGS = "create <column_name>? - optional\n"
    ERR_MSG_READ_ARGS = "read <column_name> <row>\n"
    ERR_MSG_UPDATE_ARGS = "update <column_name> <row> <value>\n"
    ERR_MSG_DELETE_ARGS = "delete <row>\n"
    ERR_MSG_FORM_REQUIRED = "A form name is required\n"
    ERR_MSG_FORM_EXISTS = "There already exists a form named '%s'.\n"
    ERR_MSG_FORM_NOT_EXISTS = "There is no form named '%s'.\n"
    ERR_MSG_FORM_create = "form-create <name> [list]\n"
    ERR_MSG_FORM_UPDATE = "form-update <name> [list]\n"
    ERR_MSG_FORM_RENAME = "form-rename <name> <src> <dest>\n"
)

//MARK - Console functions

/** print only in verbose mode */
func v(format string, args ...string) {
    if app_data.Verbose {
        fmt.Printf(format, args)
    }
}

/** print to error */
func e(format string, args ...string) {
    fmt.Fprintf(os.Stderr, format, args)
}

/** print to error but only in verbose mode */
func ev(format string, args ...string) {
    if app_data.Verbose {
        fmt.Fprintf(os.Stderr, format, args)
    }
}

/**************************************/
//MARK - Utility functions

/*
Run the external 'rpn' command
rpn -formula '2 3 +' -pop
*/
func run(formula string) string {
    ev("Calling the command: '%s'.\n", formula)
    out, err := exec.Command(app_data.Rpn, "-formula", formula, "-pop").Output()
    if err != nil {
        fmt.Printf("%s", err)
    }
    output := strings.TrimSpace(string(out[:]))
    ret := output
    return ret
}

//MARK - database functions
/**************************************/
/* helpers */

/** look for an option in an array and if not found use the fallback */
func arg (args []string, index int, fallback string) string {
    //[a, b, c, d] ; len=4
    //i==3
    ret := fallback
    if index<len(args) {    //request in range
        raw := args[index]
        if 0<len(raw) {
            ret = raw
        }
    }
    return ret
}

/** check if a string is contained in a list of strings */
func contains(arr []string, str string) bool {
   for _, a := range arr {
      if a == str {
         return true
      }
   }
   return false
}

/** return sorted keys from a map of interfaces */
func sorted_keys(data map[string][]interface{}) []string {
    keys := make([]string, len(data))
    i := 0
    for k := range data {
        keys[i] = k
        i++
    }
    sort.Strings(keys)
    return keys
}

/* util method to find the length of the 'first' column */
func data_length() int {
    return DataLength(app_data.Data)
}

/* find the length of the 'first' column */
func DataLength(data DataBase) int {
    length := -1
    for _ , v := range data.Columns {
        length = len(v)
        break
    }
    return length
}

func FirstForm(data DataBase) string{
    name := "def"
    for k, _ := range data.Forms {
        name = k
        break
    }
    return name
}

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

/** cache the calculated results */
func put_cache(key string, data []interface{}) {
    if  app_data.ColumnCache==nil {
         app_data.ColumnCache = make(map[string][]interface{})
    }
     app_data.ColumnCache[key] = data
}

/** get cached calculated results */
func get_cache(key string) []interface{} {
    if  app_data.ColumnCache==nil {
         app_data.ColumnCache = make(map[string][]interface{})
    }
    data :=  app_data.ColumnCache[key]

    return data
}

/**
convert a formula to a value
@param formula calculation to make $c1 $c2 +
@param row 0 based row count
@return result
*/
func formula_for_row(formula string, row int) string {
    return formula_for_data(formula, row,  app_data.Data)
}

func _formula_for_row(formula string, row int) string {
    words := strings.Split(formula, " ")
    for i,v := range words {
        //this allows for the row number to be inserted in as as column
        if strings.HasPrefix(v, "#row") {
            words[i] = fmt.Sprintf("%d",row)
        }
        if strings.HasPrefix(v, "$") {
            key := v[1:]
            columns :=  app_data.Data.Columns[key]
            if columns!=nil {
                column := fmt.Sprintf("%f",columns[row])
                words[i] = column
            }
        }
    }
    ret := strings.Join(words, " ")
    ret = run(ret)
    return ret
}

/**
convert a formula to a value
@param formula calculation to make $c1 $c2 +
@param row 0 based row count
@return result
*/
func formula_for_data(formula string, row int, data DataBase) string {
    words := strings.Split(formula, " ")
    for i,v := range words {
        //this allows for the row number to be inserted in as as column
        if strings.HasPrefix(v, "#row") {
            words[i] = fmt.Sprintf("%d",row)
        }
        if strings.HasPrefix(v, "$") {
            key := v[1:]
            columns := data.Columns[key]
            if columns!=nil {
                column := fmt.Sprintf("%f",columns[row])//what happens if NaN
                words[i] = column
            }
        }
    }
    ret := strings.Join(words, " ")
    ret = run(ret)
    return ret
}

/******************************************************************************/
// #mark Commands

/** Dump out a list of columns with their rows */
func List(data DataBase) {
    fmt.Printf("List: ")
    for k,v := range data.Columns {
        fmt.Printf("%s=%+v ", k, v)
    }
    for k,v := range  app_data.ColumnCache {
        fmt.Printf("%s=%+v ", k, v)
    }
    fmt.Printf("\n")
}

/* Write out the header and one row of data */
func Row(args []string, data DataBase) {
    var header bytes.Buffer
    var body bytes.Buffer

    //row form? row? delimiter?
    row := 0
    form := "main"
    delimiter := " "
    
    if 0<len(args) {
        form = arg(args, 0, FirstForm(data))
    }
    if 1<len(args) {
        //just a row number
        raw_row, err := strconv.Atoi(arg(args, 1, "0"))
        if err!=nil {
            fmt.Printf("error: %v\n", err)
            return
        } else {
            row = raw_row
        }
    }
    if 2<len(args) {
    	delimiter = arg(args, 2, " ")
    }
    keys := data.Forms[form]

    for i, v := range keys {
        if value, exists := data.Columns[v] ; exists {
            if i!=0 {
                header.WriteString(delimiter)
                body.WriteString(delimiter)
            }
            header.WriteString(v)
            body.WriteString(fmt.Sprintf("%v", value[row]))
        }
    }
    fmt.Printf("%s\n%s\n", string(header.Bytes()), string(body.Bytes()))
}

/** output the calculations by row and column */
func Calculate() {
    var header bytes.Buffer
    var rows []bytes.Buffer
    first := true

    // find the first column and get its length, then initialize rows
    for _,v := range  app_data.Data.Columns {
        for i:=0 ; i<len(v) ; i++ {
            rows = append(rows, bytes.Buffer{})
        }
        break
    }

    // calculate each formula
    for key,formula := range  app_data.Data.Calculations {
        if !first {
            header.WriteString( app_data.Format.divider )
        }
        header_title := fmt.Sprintf("%s='%v'", key, formula )
        header.WriteString( header_title )

        var calc_values []interface{}
        for i,_ := range rows {
            if !first {
                rows[i].WriteString( app_data.Format.divider )
            }
            result := formula_for_row(formula, i)

            result_as_float, err := strconv.ParseFloat(result, 64)
            if err == nil {
                calc_values = append(calc_values, result_as_float)
            } else {
                calc_values = append(calc_values, result)
            }
            rows[i].WriteString( result )
        }
        put_cache(key, calc_values)
        first = false
    }
    fmt.Printf("---\n%v\n\n", string(header.Bytes()))
    for i := range rows {
        fmt.Printf("%d: %v\n", i, string(rows[i].Bytes()))
    }
}

//MARK -

/** used by Sub only */
func value(data DataBase, form string, column int, row int) string {
    form_data :=  data.Forms[form]
    column_name := form_data[column]
    cell_data :=  data.Columns[column_name]
    value := "unknown"
    if cell_data == nil {
        value = "calc"// data.Calculation[
    } else {
    	cell := cell_data[row]
    	if is_interface_a_number(cell) {
	        value = fmt.Sprintf("%f", interface_to_float( cell ) )
    	} else if is_interface_a_string(cell) {
	        value = interface_to_string(cell)
    	} else {
    		value = "~bad~"
    	}
    }
    return value
}

//test code
func Sub(args []string, data DataBase) {
	form := arg(args, 0, "main")
    //build a grid, but how big?
    column_count := len(data.Forms[form])
  	row_count := DataLength(data)

fmt.Printf("Form: %s, size: %dx%d\n", form, column_count, row_count)

    //fill out the grid of data
    grid := make( [][]string, 0 )
    for r:=0; r<row_count; r++ {
        tmp := make( []string, 0 )
        for c:=0; c<column_count; c++ {
            tmp = append( tmp, value(data, form, c, r) )
        }
        grid = append( grid, tmp )
    }

	//find widest cells
	widths := make([]int, column_count)
    for r:=0; r<row_count; r++ {
    	for c,vv := range data.Forms[form] {
        //for c:=0; c<column_count; c++ {
	    	if r==0 {
    			widths[c] = len(vv)
    		}
        	cell_value := len(value(data, form, c, r))
        	widths[c] = int(math.Max(float64(cell_value), float64(widths[c])))
            //tmp = append( tmp,  width)
            
        }
        //append( widths, tmp )
	}

fmt.Printf("widths: %v\n", widths)

    //print out the grid
    for i,_ := range grid { //rows
        if i==0 { //first line, print header
            for ii,vv := range  data.Forms[form] {
            	cell_width := widths[ii]
                if ii==0 {
                	//format := "| %10s |"
                	format := fmt.Sprintf("| %s-%ds |", "%", cell_width)
                    fmt.Printf(format, vv)
                } else {
                	// format := " %10s |"
                	format := fmt.Sprintf(" %s-%ds |", "%", cell_width)
                    fmt.Printf(format, vv)
                }
            }
            fmt.Printf("\n")
            //print out header divider line
            for ii,_ := range  data.Forms[form] { //columns
            	cell_width := widths[ii]
                if ii==0 {
                	//format := "| %10s |"
                	format := fmt.Sprintf("| %s%ds |", "%", cell_width)
                    fmt.Printf(format, strings.Repeat("-", cell_width) )
                } else {
                	// format := " %10s |"
                	format := fmt.Sprintf(" %s%ds |", "%", cell_width)
                    fmt.Printf(format, strings.Repeat("-", cell_width) )
                }
            }
            fmt.Printf("\n")
        }
        //print out data
        for ii,vv := range grid[i] {
            cell_width := widths[ii]
            if ii==0 {
                //format := "| %10s |"
				format := fmt.Sprintf("| %s%ds |", "%", cell_width)
                fmt.Printf(format, vv)
            } else {
				// format := " %10s |"
				format := fmt.Sprintf(" %s%ds |", "%", cell_width)
                fmt.Printf(format, vv)
            }
        }
        fmt.Printf("\n")
    }
}

//MARK -

//create a sample database with 3x2 columns and rows, 2 forms, one setting
func InitDataBase() DataBase {
    data := DataBase{}

    data.Columns = make( map[string][]interface{} )
    data.Columns["foo"] = make( []interface{}, 2 )
    data.Columns["foo"] = []interface{}{0.0,1.0,2.0}
    data.Columns["bar"] = make( []interface{}, 2 )
    data.Columns["bar"] = []interface{}{3.0,4.0,3.0}
    data.Columns["rab"] = make( []interface{}, 2 )
    data.Columns["rab"] = []interface{}{5.0,6.0,6.0}

    data.Forms = make( map[string][]string )
    data.Forms["main"] = []string{"foo","bar","foobar", "row"}
    data.Forms["alt"] = []string{"bar","rab","foobar", "row"}

    data.Calculations = make ( map[string]string )
    data.Calculations["foobar"] = "$foo $bar +"
    data.Calculations["row"] = "#row"

    data.Settings = make ( map[string]string )
    data.Settings["author"] = "thomas.cherry@gmail.com"
    data.Settings["main.summary"] = "avg,sum,avg"
    data.Settings["alt.summary"] = "sum,avg,sum"

    return data
}

func Initialize(file_name string) {
    data := InitDataBase()
    fmt.Printf("the database is %+v\n", data)

    //file := "data.json"
    file := file_name
    if len(file_name)<1 {
        file =  app_data.ActiveFile
    }

    var json_text []byte
    var err error
    if  app_data.IndentFile {
        json_text, err = json.MarshalIndent(data, "", "    ")
    } else {
        json_text, err = json.Marshal(data)
    }
    if err!=nil {
        fmt.Printf("error: %s\n", err)
    }
    err = ioutil.WriteFile(file, json_text, 0644)
    if err!=nil {
        fmt.Printf("Error: %s\n", err)
    } else {
        v("File %s has been saved\n", file)
    }

}

/******************************************************************************/
// #mark - application functions

func ProcessManyLines(raw_line string, app_data *AppData) DataBase {
	data := app_data.Data
    if 0<len(raw_line) {
        commands := strings.Split(raw_line, ";")
        for _, raw_command := range commands {
            command := strings.Trim(raw_command, " ")
            if 0<len(command) {
                ProcessLine(command, app_data)
            }
        }
    }
    return data
}

//MARK Command list:

//Process a line with a command and arguments
// * @param raw line to posible execute
// * @param data database to operate on
func ProcessLine(raw string, app_data *AppData) DataBase {
	data := app_data.Data
    list := strings.Split(raw, " ")
    command := list[0]
    args := []string{""}
    if len(list)>1 {
        args = list[1:]
    }
    switch command {
        case "h", "help":
            Help()
        case "q", "quit", "exit":
            if app_data.Verbose { fmt.Printf("getting out of here\n") }
            app_data.Running = false
            os.Exit(0)
        case "e", "echo":
            fmt.Printf("%s => %s\n", command, strings.Join(args, ",") )
        case "-", "----":
            Dash(args)
        case "verbose":
            app_data.Verbose = !app_data.Verbose
            v("Verbose is %s\n", "on")
        case "file":
            if 0<len(args) && 0<len(args[0]) {
                //set mode
                 app_data.ActiveFile = args[0]
            } else {
                fmt.Printf("Active file: '%s'.\n",  app_data.ActiveFile)
            }
        case "rpn":
            if 0<len(args) && 0<len(args[0]) {
                app_data.Rpn = args[0]
            } else {
                fmt.Printf("RPN command: %s\n", app_data.Rpn)
            }

        /**************************************************************/
        /* CRUD of data */

        case "c", "create":     //create ; add row or column
            Create(args)
        case "r", "read":       //read column row
        	Read(args)
        case "u", "update":     //update column, row value
        	Update(args)
        case "d", "delete":     //delete row or column
        	Delete(args)

        case "n", "rename":     //rename a row
            src_name := arg(args, 0, "")
            dest_name := arg(args, 1, "")
            if 0<len(src_name) && 0<len(dest_name) {
                 app_data.Data.Columns[dest_name] =
                     app_data.Data.Columns[src_name]
                delete( app_data.Data.Columns, src_name)
            }

        case "a", "append":
            AppendTable(app_data.Data, args)
        case "A", "append-by-name":
            AppendTableByName(app_data.Data, args)

        /**************************************************************/
        /* Form CRUD */

        case "fc", "form-create":
            FormCreate(args)
        case "fr", "form-read":
            FormRead(args)
        case "fu", "form-update":
            FormUpdate(args)
        case "fd", "form-delete":
            FormDelete(args)
        case "fn", "form-rename":
            FormRename(args)

        /**************************************************************/
        /* Calculation CRUD */

        case "cc", "calc-create":
            CalculationCreate(args)
        case "cr", "calc-read":
            CalculationRead(args)
        case "cu", "calc-update":
            CalculationUpdate(args)
        case "cd", "calc-delete":
            CalculationDelete(args)
        case "cn", "calc-rename":
            CalculationRename (args)

        /**************************************************************/
        /* Summary CRUD */
        
        case "sc", "sum-create":
        	SummaryCreate(&app_data.Data, args)
        case "sr", "sum-read":
        	SummaryRead(app_data.Data, args)
        case "su", "sum-update":
        	SummaryUpdate(&app_data.Data, args)
        case "sd", "sum-delete":
        	SummaryDelete(&app_data.Data, args)

        /**************************************************************/
        /* Other actions */

        case "FF", "Form":
            form := ""
            action := "create"
            if 0<len(args) {
                form = args[0]
            }
            if 1<len(args) {
                action = args[1]
            }
            FormFiller(form, action) //TODO: not done
        case "ff", "form":
            form := ""
            action := "create"
            if 0<len(args) {
                form = args[0]
            }
            if 1<len(args) {
                action = args[1]
            }
            FormFillerVisual(form, action) //TODO: not done
        case "markdown?":
            fmt.Printf("markdown is %t.\n", app_data.Format.markdown)
        case "markdown":
            app_data.Format.markdown = !app_data.Format.markdown
        case "sort?":
            fmt.Printf("sort is %t.\n", app_data.Sort)
        case "sort":
            app_data.Sort = !app_data.Sort
        case "t", "table":
            Table(args[0])
        case "sum", "summary":
            form := arg(args, 0, "main")
            options := arg(args, 1,  app_data.Data.Settings[form+".summary"])
            Summary(form, options)
        case "calc", "calculate":
            Calculate() //TODO: not done
        case "init", "initialize":
            file :=  app_data.ActiveFile
            if len(args)==1 || 0<len(args[0]) {
                file = args[0]
            }
            Initialize(file)
        case "l", "ls", "list":
            List(data)
        case "row":
            Row(args, app_data.Data)
        case "-dev":
            Sub(args, app_data.Data) //- test function
        case "dump":
            DumpJson()

        /*case "cs", "calcs":
            Nop()*/
        case "s", "save":
            file :=  app_data.ActiveFile
            if len(args)==1 && 0<len(args[0]) {
                file = args[0]
            }
            Save( app_data.Data, file)
    }
    return data
}

// #mark

func SetAppData(ad AppData) {
	app_data = ad
}

func GetAppData() AppData {
	return app_data
}
