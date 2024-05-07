package lib

/*
Function for managing the CRUD of calculations in the app data base
*/

import (
	"fmt"
)

func CalculationCreate (args []string) {
    name := arg(args, 0, "")
    formula := ""
    for i:=1 ; i<len(args) ; i++ {
        formula = formula + " " + args[i]
    }
    if 0<len(name) && 0<len(formula) {
        //Calculations may be nil
        if  app_data.Data.Calculations == nil {
             app_data.Data.Calculations = make ( map[string]string )
        }
         app_data.Data.Calculations[name] = formula
    }
}

func CalculationRead (args []string) {
    name := arg(args, 0, "")
    if 0<len(name) {
        fmt.Printf("%s\n",  app_data.Data.Calculations[name])
    } else {
        fmt.Printf("%v\n",  app_data.Data.Calculations)
    }
}

func CalculationUpdate (args []string) {
    name := arg(args, 0, "")
    formula := ""
    for i:=1 ; i<len(args) ; i++ {
        formula = formula + " " + args[i]
    }
    if 0<len(name) && 0<len(formula) {
        if _, ok :=  app_data.Data.Calculations[name] ; ok {
            fmt.Printf("no calculation")
        } else {
             app_data.Data.Calculations[name] = formula
        }
    }
}

func CalculationDelete (args []string) {
    name := arg(args, 0, "")
    if 0<len(name) {
        delete( app_data.Data.Calculations, name)
    }
}

func CalculationRename (args []string) {
    src_name := arg(args, 0, "")
    dest_name := arg(args, 1, "")
    if 0<len(src_name) && 0<len(dest_name) {
         app_data.Data.Calculations[dest_name] =
             app_data.Data.Calculations[src_name]
        delete( app_data.Data.Calculations, src_name)
    }
}
