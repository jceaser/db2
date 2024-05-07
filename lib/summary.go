package lib

/*
Functions for managing the summary in the app data base
*/

import (
	"fmt"
    "strings"
    "bytes"
    "math"
    "sort"
)

const (
    SUMMARY_FUNCS = "avg, count, har, max, medium, mode, min, nop, sum, sdev"
)

/* Create a form summary in settings */
func SummaryCreate(data *DataBase, args []string) {
    if 1<len(args) {
    	form := arg(args, 0, "main")
		list := arg(args, 1, "")
		form_key := fmt.Sprintf("%s.summary", form)
		if _, okay := data.Settings[form_key]; !okay {
			data.Settings[form_key] = list
		}
	}
}

/* Create a form summary in settings */
func SummaryRead(data DataBase, args []string) {
    if 0<len(args) {
    	if 0 == len(args[0]) {
    		//no form was given
    		for key, _ := range app_data.Data.Forms {
				form_key := fmt.Sprintf("%s.summary", key)
				if list, okay := data.Settings[form_key]; okay {
					fmt.Printf("%s = %s\n", key, list)
    			}
    		}
    	} else {
    		//one form
			form := arg(args, 0, "main")
			form_key := fmt.Sprintf("%s.summary", form)
			if list, okay := data.Settings[form_key]; okay {
				fmt.Printf("%s\n", list)
			}
    	}
	}
}

/* Create a form summary in settings */
func SummaryUpdate(data *DataBase, args []string) {
    if 1<len(args) {
    	form := arg(args, 0, "main")
		list := arg(args, 1, "")
		form_key := fmt.Sprintf("%s.summary", form)
		data.Settings[form_key] = list
	}
}

/* Delete the list for a form summary in settings */
func SummaryDelete(data *DataBase, args []string) {
    if 0<len(args) {
    	form := arg(args, 0, "main")
		form_key := fmt.Sprintf("%s.summary", form)
		if _, okay := data.Settings[form_key]; okay {
			data.Settings[form_key] = ""
		}
	}
}

// Summaries a form by printing out a table, first row is header, last row is
// summary row. Each column is represented on the summary row based on data
// example: sum main avg,avg
// * @param form name of form to summarize
// * @param args dash delimitated list of summarize functions
func Summary(form string, args string) {
    var out bytes.Buffer
    if 0<len(form) {
        v("sumarize form %s with %s\n", form, args)
        if  app_data.Data.Forms[form]==nil {
            if form == "main" {
                //create_main_form()
            } else {
                fmt.Printf("Could not find form '%s'.\n", form)
                return
            }
        }
        Table(form)
        out.WriteString(" ")
        first_form :=  app_data.Data.Forms[form][0]
        var alist []string
        if len(args)<1 {
            form_summary :=  app_data.Data.Settings[form+".summary"]
            if 0<len(form_summary) {
                alist = strings.Split(form_summary, ",")
            }
        } else {
            alist = strings.Split(args, ",")
        }
        for i,value := range alist {
            if i<len( app_data.Data.Forms[form]) {
                field :=  app_data.Data.Forms[form][i]
                data :=  app_data.Data.Columns[field]
                if data == nil {
                    /*
                    there is no column data, so try getting calculated values
                    from the cache. Table caches it's last calculations for
                    functions like summary to build on
                    */
                    row_count := len( app_data.Data.Columns[first_form])
                    data = make([]interface{}, row_count)
                    raw := get_cache(field)
                    for i,cached_value := range raw {
                        data[i] = cached_value
                    }
                }
                if 0<i {
                    out.WriteString( app_data.Format.divider )
                }
                result := ""
                switch value {
                    case "a", "avg":
                        result = fmt.Sprintf(app_data.Format.template_float, Average(data) )
                    case "c", "cnt", "count":
                        result = fmt.Sprintf(app_data.Format.template_decimal, len(data))
                    case "h", "har", "harmonic":
                        result = fmt.Sprintf(app_data.Format.template_float, Harmonic(data))
                    case "mx", "max":
                        result = fmt.Sprintf(app_data.Format.template_float, Max(data))
                    case "m", "med", "medium":
                        result = fmt.Sprintf(app_data.Format.template_float, Median(data))
                    case "md", "mod", "mode":
                        result = fmt.Sprintf(app_data.Format.template_float, Mode(data))
                    case "mn", "min":
                        result = fmt.Sprintf(app_data.Format.template_float, Min(data))
                    case "n", "nop":
                        result = fmt.Sprintf(app_data.Format.template_string, "")
                    case "s", "sum":
                        result = fmt.Sprintf(app_data.Format.template_float, Sum(data))
                    case "sd", "dev", "sdev":
                        sd := StandardDeviation(data)
                        result = fmt.Sprintf(app_data.Format.template_float, sd)
                }
                out.WriteString ( result )
            }
        }
        fmt.Printf( "%v\n", string(out.Bytes()) )
    }
}

//MARK - Summary functions

func Average(data []interface{}) float64 {
    total := 0.0
    count := 0
    average := 0.0
    for _, value := range data {
        if is_interface_a_number(value) {
           total = total + interface_to_float(value)
            count = count + 1
        }
    }
    average = total / float64(count)
    return average
}

func Harmonic(data []interface{}) float64 {
    total := 0.0
    count := 0
    harmonic := 0.0
    for _, value := range data {
        if is_interface_a_number(value) {
            total = total + ( 1.0 / interface_to_float(value) )
            count = count + 1
        }
    }
    harmonic = float64(count) / total
    return harmonic
}

func StandardDeviation (data []interface{}) float64 {
    var sum, mean, sd float64 = 0, 0, 0
    count_i := 0;//len(data)
    count_f := 0.0//float64(count_i)

    for i:=0 ; i<count_i; i++ {
        if is_interface_a_number(data[i]) {
            sum += interface_to_float(data[i])
            count_i = count_i + 1
        }
    }
    count_f = float64(count_i)
    mean = sum / count_f
    for i:=0 ; i<count_i; i++ {
        if is_interface_a_number(data[i]) {
            sd += math.Pow( interface_to_float(data[i])-mean, 2)
        }
    }
    sd = math.Sqrt( sd / count_f)
    return sd
}

func Sum(data []interface{}) float64 {
    total := 0.0
    for _,value := range data {
        total = total + interface_to_float(value)
    }
    return total
}

func Max(data []interface{}) float64 {
    max := math.SmallestNonzeroFloat64
    for _,value := range data {
        if is_interface_a_number(value) {
            max = math.Max(max, interface_to_float(value))
        }
    }
    return max
}

func Min(data []interface{}) float64 {
    min := math.MaxFloat64
    for _,value := range data {
        if is_interface_a_number(value) {
            min = math.Min(min, interface_to_float(value))
        }
    }
    return min
}

func Median(data []interface{}) float64 {

    sort.Slice(data, func(i, j int) bool {
        return interface_to_float(data[i]) < interface_to_float(data[j])
    })

    len_of_data := float64(len(data))
    ret := 0.0
    if math.Mod(len_of_data, 2) == 0.0 { //even number
        index := int((math.Floor(len_of_data) / 2) - 1.0)
        left := interface_to_float(data[index])
        right := interface_to_float(data[index+1])
        ret = ( left + right ) / 2.0
    } else { //odd number
        index := int(math.Floor(len_of_data / 2))
        ret = interface_to_float(data[index])
    }
    return ret
}

func Mode(data []interface{}) float64 {
    sort.Slice(data, func(i, j int) bool {
        return interface_to_float(data[i]) < interface_to_float(data[j])
    })
    hash := make( map[float64]int )
    //collect counts
    for _, v := range data {
        value := interface_to_float(v)
        existing := hash[value]
        hash[value] = existing+1
    }
    selected := 0.0
    count := 1  //assume at least two values to bump off the default
    for k, v := range hash {
        if count<v {
            selected = k
            count = v
        }
    }
    return selected
}

/* do nothing */
func Nop() {
    fmt.Printf("not implemented yet\n")
}
