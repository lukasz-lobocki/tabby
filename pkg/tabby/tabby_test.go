/*
Prints left-aligned table.

ANSI color sequences within cells do not distort alignment.

	_tabby := new(tabby.Table)

	if err := _tabby.AddHeaders([]string{
			"Something",
			"One mo" + "\033[0;31m" + "r" + "\033[0m" + "e",
		}); err != nil {log.Fatalln(err)}

	if err := _tabby.AddRowCells([]string{
		"first",
		"seco" + "\033[0;31m" + "n" + "\033[0m" + "d_garbage67890",
	}); err != nil {log.Fatalln(err)}

	_tabby.Print(nil)
*/
package tabby

import (
	"reflect"
	"testing"
)

func Test_getRuneCount(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Pure ASCII", args{"Bulba"}, 5},
		{"UTF-8", args{"Łukasz Łobocki"}, 14},
		{"UTF-8 with ANSI colors", args{"\033[0;31mŁukasz Ł\033[0mobocki"}, 14},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRuneCount(tt.args.input); got != tt.want {
				t.Errorf("getRuneCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_padRight(t *testing.T) {
	type args struct {
		input   string
		lenght  int
		padding string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Pure ASCII", args{"Bulba", 12, "."}, "Bulba......."},
		{"UTF-8 with ANSI colors", args{"\033[0;31mŁukasz Ł\033[0mobocki", 20, "."}, "\033[0;31mŁukasz Ł\033[0mobocki......"},
		{"UTF-8 with ANSI colors exact lenght", args{"\033[0;31mŁukasz Ł\033[0mobocki", 14, "."}, "\033[0;31mŁukasz Ł\033[0mobocki"},
		{"UTF-8 with ANSI colors exeeds lenght", args{"\033[0;31mŁukasz Ł\033[0mobocki", 1, "."}, "\033[0;31mŁukasz Ł\033[0mobocki"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := padRight(tt.args.input, tt.args.lenght, tt.args.padding); got != tt.want {
				t.Errorf("padRight() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getDefaultConfig(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		{"Single", &Config{" ", "  "}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDefaultConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDefaultConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_formatTableLine(t *testing.T) {
	type args struct {
		_l            []string
		_columnsWidth []int
		padding       string
		spacing       string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Empty", args{}, ""},
		{"Simple", args{_l: []string{"a", "b"}, _columnsWidth: []int{2, 2}, padding: ".", spacing: "|"}, "a.|b."},
		{"Overflowing", args{_l: []string{"abc", "b"}, _columnsWidth: []int{2, 2}, padding: ".", spacing: "|"}, "abc|b."},
		{"Longer", args{_l: []string{"a", "b"}, _columnsWidth: []int{3, 2}, padding: ".", spacing: "|"}, "a..|b."},
		{"UTF-8", args{
			_l:            []string{"\033[0;31mŁukasz Ł\033[0mobocki12345", "b"},
			_columnsWidth: []int{21, 2},
			padding:       "+",
			spacing:       "|"},
			"\033[0;31mŁukasz Ł\033[0mobocki12345++|b+"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatTableLine(tt.args._l, tt.args._columnsWidth, tt.args.padding, tt.args.spacing); got != tt.want {
				t.Errorf("formatTableLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getColumnsWidth(t *testing.T) {
	type args struct {
		_t Table
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{"Empty", args{}, []int{}},
		{"Just single header", args{Table{Header{"a"}, rows{}}}, []int{1}},
		{"Just double header", args{Table{Header{"a", "\033[0;31mŁukasz Ł\033[0mobocki"}, rows{}}}, []int{1, 14}},
		{"UTF-8", args{Table{
			Header{"a", "\033[0;31mŁukasz Ł\033[0mobocki"},
			rows{Row{
				"\033[0;31mŁukasz Ł\033[0mobocki12345",
				"b",
			}, Row{"c", "d"}}}}, []int{19, 14}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getColumnsWidth(tt.args._t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getColumnsWidth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_Print(t *testing.T) {
	type fields struct {
		header Header
		rows   rows
	}
	type args struct {
		config *Config
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Empty", fields{}, args{nil}},
		{"Populated", fields{header: Header{"A"}, rows: rows{Row{"1"}, Row{"2"}}}, args{nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_t := &Table{
				header: tt.fields.header,
				rows:   tt.fields.rows,
			}
			_t.Print(tt.args.config)
		})
	}
}

func TestTable_SetHeader(t *testing.T) {
	type fields struct {
		header Header
		rows   rows
	}
	type args struct {
		header Header
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Empty", fields{}, args{}, true},
		{"Single header", fields{}, args{header: Header{"first"}}, false},
		{"Double header", fields{}, args{header: Header{"first", "second"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_t := &Table{
				header: tt.fields.header,
				rows:   tt.fields.rows,
			}
			if err := _t.SetHeader(tt.args.header); (err != nil) != tt.wantErr {
				t.Errorf("Table.SetHeader() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTable_AppendRow(t *testing.T) {
	type fields struct {
		header Header
		rows   rows
	}
	type args struct {
		row Row
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Empty", fields{}, args{}, false},
		{"Single header", fields{header: Header{"one"}}, args{row: Row{"first"}}, false},
		{"Double header", fields{header: Header{"one", "two"}}, args{row: Row{"first", "second"}}, false},
		{"More cells than headers", fields{header: Header{"one"}}, args{row: Row{"first", "second"}}, true},
		{"Less cells than headers", fields{header: Header{"one", "two"}}, args{row: Row{"first"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_t := &Table{
				header: tt.fields.header,
				rows:   tt.fields.rows,
			}
			if err := _t.AppendRow(tt.args.row); (err != nil) != tt.wantErr {
				t.Errorf("Table.AppendRow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
