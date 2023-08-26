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
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Table contents to be emitted.
type Table struct {
	headers []string   // slice of headers
	rows    [][]string // slice of slices of cells
}

// Strings used for padding and spacing.
type Config struct {
	padding string // chars added to the right of cell contents
	spacing string // chars added between cells
}

/*
Adds headers to the table.

	params: slice of headers

Headers must be added before adding rows.
*/
func (_t *Table) AddHeaders(headers []string) error {

	// Error if no headers provided
	if len(headers) < 1 {
		return errors.New(fmt.Sprint("no headers provided."))
	}

	*_t = Table{
		headers: headers,
	}
	return nil
}

/*
Appends row of cells to the table.

	params: slice of cells

Headers must be added before adding rows. Number of cells must not exceed number of headers.
*/
func (_t *Table) AddRowCells(row []string) error {

	// Error if number of cells in the row exceeds the number of headers
	if len(row) > len(_t.headers) {
		return errors.New(
			fmt.Sprintf("number of cells %d in the row %v exceeds the number of headers %d.",
				len(row),
				row,
				len(_t.headers)))
	}

	_t.rows = append(_t.rows, row)
	return nil
}

/*
Prints the table.

	params: config structure (optional)
*/
func (_t *Table) Print(config *Config) {

	if config == nil {
		// defaultConfig returns the default config for table
		config = getDefaultConfig()
	}

	// Measure columns for biggest widht
	_columnsWidth := getColumnsWidth(*_t)

	// Emit header
	fmt.Println(
		formatTableLine(
			_t.headers,
			_columnsWidth,
			config.padding,
			config.spacing))

	// Iterate and emit rows
	for _, _row := range _t.rows {
		// Emit row
		fmt.Println(
			formatTableLine(
				_row,
				_columnsWidth,
				config.padding,
				config.spacing))
	}
	return
}

// Provides the default config for table
func getDefaultConfig() *Config {
	return &Config{
		padding: " ",
		spacing: "  ",
	}
}

// Formats table line appending cells, padding to given width and spacing between the cells
func formatTableLine(_l []string, _columnsWidth []int, padding string, spacing string) string {

	var _ln strings.Builder

	for i, _cell := range _l {
		// Append each cell padded
		_ln.WriteString(padRight(_cell, _columnsWidth[i], padding))
		if i < len(_l)-1 {
			// Append spacing but not after last column; subtracting -1 to adjust 0-based loop
			_ln.WriteString(spacing)
		}
	}
	return _ln.String()
}

// Returns string padded to visible rune lenght.
func padRight(input string, lenght int, padding string) string {

	// If input not shorter than lenght, return input
	if _runeCount := getRuneCount(input); _runeCount >= lenght {
		return input
	}

	return input + strings.Repeat(padding, lenght-getRuneCount(input))
}

// Returns longest runic lenght of each column with header.
func getColumnsWidth(_t Table) []int {

	_output := make([]int, len(_t.headers))

	// Measure header
	for i, _header := range _t.headers {
		_output[i] = getRuneCount(_header)
	}

	// Iterate and measure rows
	for _, _row := range _t.rows {
		for j, _cell := range _row {
			if _thisLength := getRuneCount(_cell); _thisLength > _output[j] {
				_output[j] = _thisLength
			}
		}
	}
	return _output
}

// Returns string with ANSI codes removed.
func removeANSICodes(input string) string {
	_regexp := regexp.MustCompile(`\x1b\[[0-9;]*[mK]`)
	return _regexp.ReplaceAllString(input, "")
}

// Returns rune count with ANSI codes removed.
func getRuneCount(input string) int {
	return utf8.RuneCountInString(removeANSICodes(input))
}
