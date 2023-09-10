/*
Prints left-aligned table. ANSI color sequences within cells do not distort the alignment.

Example:

	table := new(tabby.Table)

	if err := table.SetHeader(tabby.Header{
		"\033[4mFIRST\033[0m",
		"\033[4mSECOND\033[0m",
	}); err != nil {
		log.Fatalln(err)
	}

	if err := table.AppendRow(tabby.Row{
		"eins \033[4;33mzwei\033[0m drei",
		"vier",
	}); err != nil {
		log.Fatalln(err)
	}

	table.Print(nil)
*/
package tabby

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Elements of a table: Header and []Row.
type (
	line   []string // slice of cells (strings).
	Header line     // header line, slice of header cells, []string.
	Row    line     // row line, slice of row cells, []string.
	rows   []Row    // slice of rows.
)

// Contents of a table.
type Table struct {
	header Header // header of a table.
	rows   rows   // rows of a table.
}

// Strings used for padding and spacing.
type Config struct {
	Padding string // chars added to the right of cell contents, defaults to " ".
	Spacing string // chars added between cells, defaults to "  ".
}

/*
Sets header of the table.

'header': slice of header cells.

Header must be added before adding rows.
*/
func (_t *Table) SetHeader(header Header) error {

	// Error if no headers provided
	if len(header) < 1 {
		return errors.New(fmt.Sprint("no header provided."))
	}

	_t.header = header
	return nil
}

/*
Appends row of cells to the table.

'row': slice of row cells

Header must be set before appending rows. Number of cells in the row must not exceed number of cells in the header.
*/
func (_t *Table) AppendRow(row Row) error {

	// Error if number of cells in the row exceeds the number of headers
	if len(row) > len(_t.header) {
		return errors.New(
			fmt.Sprintf("number of cells %d in the row %v exceeds the number of cells in the header %d.",
				len(row),
				row,
				len(_t.header)))
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
			line(_t.header),
			_columnsWidth,
			config.Padding,
			config.Spacing))

	// Iterate and emit rows
	for _, _row := range _t.rows {
		// Emit row
		fmt.Println(
			formatTableLine(
				line(_row),
				_columnsWidth,
				config.Padding,
				config.Spacing))
	}
	return
}

/*
Formats table line appending cells, padding to given width and spacing between the cells.
*/
func formatTableLine(_l line, _columnsWidth []int, padding string, spacing string) string {

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

/*
Returns string padded to visible rune lenght.
*/
func padRight(input string, lenght int, padding string) string {

	// If input not shorter than lenght, return input
	if _runeCount := getRuneCount(input); _runeCount >= lenght {
		return input
	}

	return input + strings.Repeat(padding, lenght-getRuneCount(input))
}

/*
Returns longest runic lenght of each column with header.
*/
func getColumnsWidth(_t Table) []int {

	_output := make([]int, len(_t.header))

	// Measure header
	for i, _header := range _t.header {
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

/*
Provides the default config for table.
*/
func getDefaultConfig() *Config {
	return &Config{
		Padding: " ",
		Spacing: "  ",
	}
}

/*
Returns rune count with ANSI codes removed.
*/
func getRuneCount(input string) int {
	return utf8.RuneCountInString(removeANSICodes(input))
}

/*
Returns string with ANSI codes removed.
*/
func removeANSICodes(input string) string {
	_regexp := regexp.MustCompile(`\x1b\[[0-9;]*[mK]`)
	return _regexp.ReplaceAllString(input, "")
}
