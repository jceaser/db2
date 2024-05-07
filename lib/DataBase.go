package lib

import (
	//"fmt"
)

var (
	ExampleDataBaseJson = []byte(`{
	"Columns": {
		"a": [0, 3],
		"b": [1, 4],
		"c": [2, 5]
	},
	"Forms": {"alt": ["a", "c"]}
}`)
)

type DataBase struct {
    Columns map[string][]interface{}
    Forms map[string][]string
    Calculations map[string]string
    Settings map[string]string
}

func (self *DataBase) CleanUp() {
	self.EnsureMainForm()
	self.EnsureSettings()
}

func (self *DataBase) EnsureMainForm() {
    if self.Forms["main"] == nil {
        keys := sorted_keys(self.Columns)
        if  self.Forms == nil {
             self.Forms = make( map[string][]string )
        }
         self.Forms["main"] = keys
    }
}

func (self *DataBase) EnsureSettings() {
	if self.Settings == nil {
		self.Settings = map[string]string{}
	}
}

/*
func init() {
	fmt.Printf("---\n")

	db, err := JsonToStruct[DataBase](ExampleDataBaseJson)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v\n", db)
	}
	
	db.EnsureMainForm()
	
	json, err := StructToJson[DataBase](db, true)
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v\n", string(json))
	}

	fmt.Printf("---\n")
}
*/
