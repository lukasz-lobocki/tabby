package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/lukasz-lobocki/tabby/pkg/utils"
)

type Table struct {
	headers []string
	rows    [][]string
}

type Config struct {
	padding string
	spacing string
}

func main() {

	_tab := new(Table)

	if err := _tab.AddHeaders([]string{
		"something",
		"bnother",
		"one mo" + utils.RED + "r" + utils.RESET + "e",
	}); err != nil {
		log.Fatalln(err)
	}

	if err := _tab.AddRowCells([]string{
		"uno",
		"dos",
		"tres",
		//"quatro",
	}); err != nil {
		log.Fatalln(err)
	}

	if err := _tab.AddRowCells([]string{
		"jeden",
		"kl" + utils.RED + "m" + utils.RESET + "no67890",
		"trzy",
	}); err != nil {
		log.Fatalln(err)
	}

	_tab.Print(nil)
}

// Adds headers to the table
func (_t *Table) AddHeaders(headers []string) error {

	// Error if no headers provided
	if len(headers) < 1 {
		return errors.New(fmt.Sprint("no headers provided"))
	}

	*_t = Table{
		headers: headers,
	}
	return nil
}

// Adds row of cells to the table
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

// Prints the table
func (_t Table) Print(c *Config) {

	if c == nil {
		// defaultConfig returns the default config for table
		c = getDefaultConfig()
	}

	// Measure columns for biggest widht
	_columnsWidth := getColumnsWidth(_t)

	// Emit header
	fmt.Println(
		formatTableLine(
			_t.headers,
			_columnsWidth,
			c.padding,
			c.spacing))

	// Iterate and emit rows
	for _, _row := range _t.rows {
		// Emit row
		fmt.Println(
			formatTableLine(
				_row,
				_columnsWidth,
				c.padding,
				c.spacing))
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
			// Append spacing but not after last column; adding +1 to adjust 0-based loop
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
