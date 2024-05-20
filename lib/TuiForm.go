package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Contact struct {
	Title     string
	FirstName string
	LastName  string
	Address   string
	Notes     string
	Adult     bool
	Password  string
}

func (self Contact) json() []byte {
	b, err := json.Marshal(self)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return b
}

func setup_log() *os.File {
	// open a file
	f, err := os.OpenFile("test.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	//defer f.Close()
	// assign it to the standard logger
	log.SetOutput(f)
	return f
}

func file_write(file_name string, data []byte) {
	// open output file
	fo, err := os.Create(file_name)
	if err != nil {
		panic(err)
	}
	defer fo.Close()
	if _, err := fo.Write(data); err != nil {
		panic(err)
	}
}

func _save(app *tview.Application, form *tview.Form) {
	contact := Contact{}
	_, title_text := form.GetFormItemByLabel("Title").(*tview.DropDown).GetCurrentOption()

	contact.Title = title_text
	contact.FirstName = form.GetFormItemByLabel("First name").(*tview.InputField).GetText()
	contact.LastName = form.GetFormItemByLabel("Last name").(*tview.InputField).GetText()
	contact.Address = form.GetFormItemByLabel("Address").(*tview.TextArea).GetText()
	contact.Notes = form.GetFormItemByLabel("Notes").(*tview.TextView).GetText(true)
	contact.Adult = form.GetFormItemByLabel("Age 18+").(*tview.Checkbox).IsChecked()
	contact.Password = form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
	log.Printf("json: %s\n", string(contact.json()))
	file_write("contact.json", contact.json())
}

func save(app *tview.Application, form_view *tview.Form, form string, data *DataBase) {
	CreateRow(*data)
	for _, column := range data.Forms[form] {
		if 0 < len(app_data.Data.Calculations[column]) {
			continue // this is a calculation, skip it
		} else {
			index := len(data.Columns[column]) - 1
			text := form_view.GetFormItemByLabel(column).(*tview.InputField).GetText()
			text = strings.TrimSpace(text)
			text = text[1:]
			data.Columns[column][index] = tview.Escape(text)
		}
	}
}

func CreateTable(form string, data DataBase) *tview.Table {
	table := tview.NewTable().SetBorders(true)
	PopulateRows(form, data, table)
	return table
}

func PopulateRows(form string, data DataBase, table *tview.Table) {
	for c, v := range data.Forms[form] {
		//write header
		header_cell := tview.NewTableCell(v).
			SetAlign(tview.AlignCenter).
			SetTextColor(tcell.ColorRed).
			SetStyle(tcell.Style{}.Bold(true))
		table.SetCell(0, c, header_cell)
		//write rows
		for r, raw := range data.Columns[v] {
			value := interface_to_string(raw)
			table.SetCell(r+1, c, tview.NewTableCell(value))
		}
	}
}

func FormFiller(form string, action string, data *DataBase) {
	log_file := setup_log()
	defer log_file.Close()

	if form == "" {
		form = "main"
	}

	app := tview.NewApplication()

	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	form_view := tview.NewForm()
	table_view := CreateTable(form, *data)

	flex.AddItem(form_view, 0, 1, true)
	flex.AddItem(table_view, 0, 1, false)

	//set up form
	form_view.AddTextView("Form", form, 40, 2, true, false)

	for i, v := range data.Forms[form] {
		form_view.AddInputField(v, string(i), 20, nil, nil)
	}

	//form_view.AddCheckbox("Age 18+", false, nil).
	form_view.AddButton("Save", func() {
		save(app, form_view, form, data)
		//update the table view
		PopulateRows(form, *data, table_view)
	})
	form_view.AddButton("Quit", func() {
		app.Stop()
	})
	form_view.SetBorder(true).SetTitle("Enter some data").SetTitleAlign(tview.AlignLeft)
	if err := app.SetRoot(flex, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
		panic(err)
	}
}
