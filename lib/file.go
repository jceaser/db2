package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

/**************************************/
/* Manage load and unload database functions */
//MARK - Data file functions

/* set the internal data value, use this to setup tests */
func SetData(data DataBase) {
	app_data.Data = data
}

/* Takes a byte array of some string, and converts that to a Target */
func JsonToStruct[Target any](bytes []byte) (Target, error) {
	var data Target
	err := json.Unmarshal(bytes, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

/* Convent a structure to a json byte array */
func StructToJson[T any](data T, useIndent bool) ([]byte, error) {
	var bytes []byte
	var err error
	if useIndent {
		bytes, err = json.MarshalIndent(data, "", strings.Repeat(" ", 4))
	} else {
		bytes, err = json.Marshal(data)
	}
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

/* Load a json data file */
func load_json(file_path string) *os.File {
	json_raw, err := os.Open(file_path)
	if err != nil {
		if os.IsNotExist(err) {
			//create the file because it does not exist
			v("Creating data file %s\n", file_path)
			// TODO: this looks like the wrong default
			sample := []byte("{}")
			err := ioutil.WriteFile(file_path, sample, 0644)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return nil
			}
			//try to open it a second time
			json_raw, err = os.Open(file_path)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return nil
			}
		} else {
			fmt.Printf("Error: %s\n", err)
			return nil
		}
	}
	//defer json_raw.Close()
	return json_raw
}

/* Load a database from a file */
func Load(file_path string) *DataBase {
	v("Loading file %s\n", file_path)
	json_raw := load_json(file_path)
	if json_raw == nil {
		fmt.Printf("No data\n")
	} else {
		defer json_raw.Close()
		var json_data = DataBase{}
		bytes, err := ioutil.ReadAll(json_raw)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		} else {
			json.Unmarshal([]byte(bytes), &json_data)
			json_data.CleanUp()
			return &json_data
		}
	}
	return nil
}

/** save the database to a file */
func Save(data DataBase, file string) {
	var json_text []byte
	var err error
	if app_data.IndentFile {
		json_text, err = json.MarshalIndent(data, "", "    ")
	} else {
		json_text, err = json.Marshal(data)
	}

	if len(file) < 1 {
		file = app_data.ActiveFile
	}

	if err != nil {
		fmt.Printf("error: %s - %s\n", file, err)
		return
	}
	err = ioutil.WriteFile(file, json_text, 0644)
	if err != nil {
		fmt.Printf("Error: %s - %s\n", file, err)
	} else {
		v("File %s has been saved\n", file)
	}
}

/** print out json */
func DumpJson() {
	fmt.Println("Col data: ", app_data.Data.Columns)
	var json_text []byte
	var err error
	if app_data.IndentFile {
		json_text, err = json.MarshalIndent(app_data.Data, "", "    ")
	} else {
		json_text, err = json.Marshal(app_data.Data)
	}
	if err == nil {
		fmt.Printf("%s\n", json_text)
	} else {
		fmt.Printf("Error: %s\n", err)
	}
}
