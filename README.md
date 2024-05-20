# db2
A command line database built around JSON that supports forms, inspired by old versions of File Maker Pro

Name db2 is because this tool was considered a more complicated iteration of the db command found in the https://github.com/jceaser/gotools/

## About

db2 manages three main types of data in the database: Columns, Calculations, Forms, and Summaries.

* Columns, primary data. Columns contain rows of cells. Cells can contain strings, numbers, booleans.
* Calculations, RPN formulas. Columns are referenced by name with a preceding '$'
* Forms are collections of columns and calculations
* Summaries are lists of functions that operate on form columns.

There are commands to operate on all of these data types.

## Usage

db2 is a command, launch it with `db2`. Once launched commands can be issued at the prompt.

### Flags

No flag is required.

| Flag     | Argument | Description |
| -------- | -------- | ----------- |
| -command | string   | Run one command (string) and **exit**.
| -file    | string   | Data file (default "data.json").
| -init    | string   | Run initial command (string) and stay open.
| -rpn     | string   | Command to process calculations (default "rpn").
| -verbose |          | Verbose mode.

### Commands

To get a list of commands, use `help`. It will print out a list like below.

Note: Arguments with ? are optional

| Short | Long      | Arguments           | Description |
| ----- | ---------- | ------------------- | ----------- |
|  c    | create     | name         | Create a column by name
|       |            |              | create a zero row
|  r    | read       | col row      | read a column row
|  u    | update     | col row val  | update a column row
|  d    | delete     | index        | delete a row by number
|       |            | name         | delete a column by name
|  a    | append     | <value list> | append a table
|  A    | append-by-name | name:value... | append a table with named columns
|  n    | rename     | src dest     | rename a column from src to dest
|
|  fc   | form-create | name list    | Create a form
|  fr   | form-read   | name?        | list forms, all if name is not given
|  fu   | form-update | name formula | update a form
|  fd   | form-delete | name         | delete a form
|  fn   | form-rename | src dest     | rename a form from src to dest
|  ff   | form        | name action? | Create a row using a form
|  FF   | Form        | name action? | Basic form filler
|
|  cc   | calc-create | name formula | Create a calculation
|  cr   | calc-read   | name?        | list calculations, all if name is not given
|  cu   | calc-update | name formula | update a calculation
|  cd   | calc-delete | name         | delete a calculation
|  cn   | calc-rename | src dest     | rename a calculation from src to dest
|
|  sc   | sum-create  | form list    | assign a list of sum functions
|  sr   | sum-read    | form?        | read all values, or one sum function
|  su   | sum-update  | form list    | update a list of sum functions
|  sd   | sum-delete  | form         | clear a list of sum functions
|
| sum   | summary     | form list    | summarize a form with function list:
|       |             |              | avg,count,max,min,medium,mode,min,nop,sum,sdev
|
|   t   | table       | form?        | display a table, optionally as a form
|   l   | ls list     |              |
|       | row         | form row?    | show a header and one row from form
|
|   s   | save        |              | save database to file
|       | dump        |              | output the current data
|   q   | quit        |              | quit interactive mode
|       | exit        |              | quit interactive mode
|   h   | help        |              | this output
|   e   | echo        | string       | echo out something
|   -   | ----        | sep count    | print out a separator
|       | file        | name?        | set or print current file name
|       | rpn         | path?        | set or print current rpn command
|       | verbose     |              | toggle verbose mode
|       | sort?       |              | output current sorting state
|       | sort        |              | toggle the current sort mode

Each of the 4 types of data have Create Read Update Delete functions and some support more functions. These are grouped together in the help output and have similar names.

There are also functions for displaying data, like `table` and `summary`. Input can be done with the CRUD functions for columns or by using the form-filler functions.
