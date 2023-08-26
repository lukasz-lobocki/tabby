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

func TestTable_AddHeaders(t *testing.T) {
	type fields struct {
		headers []string
		rows    [][]string
	}
	type args struct {
		headers []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Empty", fields{}, args{}, true},
		{"Single header", fields{}, args{[]string{"first"}}, false},
		{"Two headers", fields{}, args{[]string{"first", "second"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_t := &Table{
				headers: tt.fields.headers,
				rows:    tt.fields.rows,
			}
			if err := _t.AddHeaders(tt.args.headers); (err != nil) != tt.wantErr {
				t.Errorf("Table.AddHeaders() error = %v, wantErr %v", err, tt.wantErr)
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

func TestTable_AddRowCells(t *testing.T) {
	type fields struct {
		headers []string
		rows    [][]string
	}
	type args struct {
		row []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Empty", fields{}, args{}, false},
		{"Single header", fields{headers: []string{"one"}}, args{[]string{"first"}}, false},
		{"Two headers", fields{headers: []string{"one", "two"}}, args{[]string{"first", "second"}}, false},
		{"More cells than headers", fields{headers: []string{"one"}}, args{[]string{"first", "second"}}, true},
		{"Less cells than headers", fields{headers: []string{"one", "two"}}, args{[]string{"first"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_t := &Table{
				headers: tt.fields.headers,
				rows:    tt.fields.rows,
			}
			if err := _t.AddRowCells(tt.args.row); (err != nil) != tt.wantErr {
				t.Errorf("Table.AddRowCells() error = %v, wantErr %v", err, tt.wantErr)
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
		{"Just single header", args{Table{headers: []string{"a"}}}, []int{1}},
		{"Just double header", args{Table{headers: []string{"a", "\033[0;31mŁukasz Ł\033[0mobocki"}}}, []int{1, 14}},
		{"saa", args{Table{
			headers: []string{"a", "\033[0;31mŁukasz Ł\033[0mobocki"},
			rows:    [][]string{{"\033[0;31mŁukasz Ł\033[0mobocki12345", "b"}, {"c", "d"}}}},
			[]int{19, 14}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getColumnsWidth(tt.args._t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getColumnsWidth() = %v, want %v", got, tt.want)
			}
		})
	}
}
