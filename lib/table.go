package lib

/*
Functions for managing tables in the app database
*/

import (
	"fmt"
	"bytes"
	"strconv"
	"strings"
	"math"
	"sort"
)

/*
outputs a number of characters to visually separate out the output
@param arg 1 if empty, output '----'
@param arg 1 if not empty, output that character
@param arg 2 if number, output arg[1] this many times
*/
func Dash(args []string) {
    if len(args)<1 {
        fmt.Printf("----\n")
    } else if len(args)==1 {
        if 0==len(args[0]) { args[0] = "----" }
        fmt.Printf("%s\n", args[0])
    } else {
        letter := args[0]
        count, err := strconv.Atoi(args[1])
        if err==nil {
            fmt.Printf("%s\n", strings.Repeat(letter, count) )
        } else {
            fmt.Printf("%s\n", letter)
        }
    }
}

/**
level - 0=top, 1=middle, 2=bottom
*/
func table_divider (markdown bool, columns int, level int) string {
    var sbuf bytes.Buffer
    if !markdown {
        if level==0 {
            sbuf.WriteRune(RuneULCorner) //'┌'
        } else if level==1 {
            sbuf.WriteRune(RuneLTee) //'├'
        } else if level==2 {
            sbuf.WriteRune(RuneLLCorner) //'└'
        }
    }
    for i:=0 ; i<=columns ; i++ {
        if markdown {
            sbuf.WriteString(app_data.Format.divider_pipe)
            //width specific
            sbuf.WriteString(fmt.Sprintf(app_data.Format.template_string, "---------"))
        } else {
            if 0<i && i<=columns {
                if level==0 {
                    sbuf.WriteRune(RuneTTee) //'┬'
                } else if level==1 {
                    sbuf.WriteRune(RunePlus) //"┼"
                } else if level==2 {
                    sbuf.WriteRune(RuneBTee) //┴
                }
            }
            sbuf.WriteString(fmt.Sprintf(app_data.Format.template_string, "──────────"))
        }
    }
    if markdown {
        sbuf.WriteString(app_data.Format.divider_pipe)
    } else {
        if level==0 {
            sbuf.WriteRune(RuneURCorner)
        } else if level==1 {
            sbuf.WriteRune(RuneRTee)
        } else if level==2 {
            sbuf.WriteRune(RuneLRCorner)
        }
    }
    return sbuf.String()
}

//Dump table of all columns
//* @param form name of the form to dump out, empty for entire table
func Table(form string) {
    markdown := app_data.Format.markdown
    var divider string
    if markdown {
        divider = app_data.Format.divider_pipe
    } else {
        divider = app_data.Format.divider
    }
    header, rows, keys := table_worker(form, divider)

    //print out the table
    if !markdown {
        fmt.Printf("%s\n", table_divider(markdown, len(keys), 0))
    }
    fmt_lined := (app_data.Format.template_string + "\n")
    fmt.Printf( fmt_lined , string(header.Bytes()))
    fmt.Printf("%s\n", table_divider(markdown, len(keys), 1))
    for i := range rows {
        fmt.Printf("%v\n", string(rows[i].Bytes()))
    }
    if !markdown {
        fmt.Printf("%s\n", table_divider(markdown, len(keys), 2))
    }
}

func table_worker(form string, divider string) (bytes.Buffer, []bytes.Buffer, []string) {
    var header bytes.Buffer
    var rows []bytes.Buffer
    var keys []string

    for i:=0 ; i<data_length() ; i++ {
       rows = append(rows, bytes.Buffer{})
    }

    first := true

    //figure out which fields need to be displayed
    if 0<len(form) {
        //use the form list
        keys =  app_data.Data.Forms[form]
    } else {
        //always sort because map is not order consistent
        keys = sorted_keys( app_data.Data.Columns)
    }

    //loop throug all the column and calculation keys
    max := len(keys)-1
    for index,k := range keys {
        last := false
        if max<=index {
            last = true
        }
        var formula = ""
        values :=  app_data.Data.Columns[k] //return a list of strings

        // if values is nil, then not a column, search calculations
        if values==nil {
            formula =  app_data.Data.Calculations[k]
            if formula=="" {
                continue //key is blank, skip it
            }
            var calc_values []interface{}
            for i,_ := range rows {
                result := formula_for_row(formula, i)
                result_as_float, _ := strconv.ParseFloat(result, 64)
                calc_values = append(calc_values, result_as_float)
            }
            put_cache(k, calc_values)
            values = calc_values
        }

        //find widest value
        max_width := 0
        {
            max_width_f := 0.0
            for i := range values {
                txt := fmt.Sprintf("%v", values[i])
                max_width_f = math.Max(float64(len(txt)), max_width_f)
            }
            max_width = int(max_width_f)
        }

        if first {

        }


        header.WriteString( divider )

        head_format := app_data.Format.template_string
        //head_format = strings.Replace(head_format, "%10s", fmt.Sprintf("%%%ds", max_width+2), 1)
        header.WriteString( fmt.Sprintf(head_format, k) )

        if last {
            //format := "%s" + fmt.Sprintf("%%%ds", max_width + 2) + "%s"
            //fmt.Printf("max=%d => %s\n", max_width, format)
            max_width = max_width
            format := "%s%10s%s"
            header.WriteString(fmt.Sprintf(format, divider, "row", divider ))
        }
        for i := range values {
            rows[i].WriteString( divider )
            column := ""
            if i<len(values) {
                format := app_data.Format.template_float
                if is_interface_a_string(values[i]) {
                    format = app_data.Format.template_string
                    /*format = strings.Replace(format,
                        "%10s",
                        fmt.Sprintf("%%%ds", max_width+2), 1)*/
                }
                column = fmt.Sprintf(format, values[i])
            }
            rows[i].WriteString( column )
            if last {
                rows[i].WriteString(fmt.Sprintf("%s%10d%s", divider, i, divider))
            }
        }
        first = false
    }

    if app_data.Sort {
        sort.Slice(rows, func(i, j int) bool {
            return rows[i].String() < rows[j].String()
        })
    }
    return header, rows, keys
}


/** append a new row to the data, and populate the named rows */
func AppendTableByName(data DataBase, args []string) {
    /* format: column_values */
    arg_count := len(args)
    if arg_count<1 {
        return
    }
    CreateRow(data)
    row := DataLength(data)-1
    for i:=0; i<len(args); i++ {
        raw := args[i]
        parts := strings.Split(raw, ":")
        column := parts[0]
        value := parts[1]
        if _, ok := data.Columns[column]; ok {
            data.Columns[column][row] = value
        }
    }
}

/** populate a new row with provided data */
func AppendTable(data DataBase, args []string) {
    /* format: column_values */
    arg_count := len(args)
    column_count := len(data.Columns)
    if arg_count<1 {
        return
    }
    CreateRow(data)
    row := DataLength(data)-1
    index := 0
    for _, column := range sorted_keys(data.Columns) {
        if value, err := strconv.ParseFloat( args[index], 64 ) ; err==nil {
            data.Columns[column][row] = value
        } else {
            data.Columns[column][row] = args[index]
        }

        //prep for next round
        index = index + 1
        if arg_count<=index || column_count<=index {
            break
        }
    }
}
