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

type table struct {
	Headers []string
	Rows    [][]string
}

func main() {

	_tab := new(table)

	if err := _tab.addHeaders([]string{
		"something",
		"bnother",
		"one mo" + utils.RED + "r" + utils.RESET + "e",
	}); err != nil {
		log.Fatalln(err)
	}

	if err := _tab.addRow([]string{
		"uno",
		"dos",
		"tres",
		//"quatro",
	}); err != nil {
		log.Fatalln(err)
	}

	if err := _tab.addRow([]string{
		"jeden",
		"kl" + utils.RED + "m" + utils.RESET + "no67890",
		"trzy",
	}); err != nil {
		log.Fatalln(err)
	}

	if err := _tab.Print(
		"  ",
		" ",
	); err != nil {
		log.Fatalln(err)
	}

}

func (_t *table) addHeaders(headers []string) error {

	// Error if no headers provided
	if len(headers) < 1 {
		return errors.New(fmt.Sprint("no headers provided"))
	}

	*_t = table{
		Headers: headers,
	}
	return nil
}

func (_t *table) addRow(row []string) error {

	// Error if number of cells in the row exceeds the number of headers
	if len(row) > len(_t.Headers) {
		return errors.New(
			fmt.Sprintf("number of columns %d in the row [%s] exceeds the number of headers %d.",
				len(row),
				row[0],
				len(_t.Headers)))
	}

	_t.Rows = append(_t.Rows, row)
	return nil
}

// Prints the table
func (_t table) Print(spacing string, padding string) error {

	const (
		DEFAULT_PADDING string = " "  // Pad with single space
		DEFAULT_SPACING string = "  " // Space with 2 spaces
	)

	// Measure columns for biggest widht
	_columnsWidth, err := getColumnsWidth(_t)
	if err != nil {
		return fmt.Errorf("error measuring column width. %w", err)
	}

	// Default padding
	if getRuneCount(padding) != 1 {
		padding = DEFAULT_PADDING
	}

	// Default spacing
	if getRuneCount(spacing) == 0 {
		spacing = DEFAULT_SPACING
	}

	// Emit header
	fmt.Println(
		formatTableLine(
			_t.Headers,
			_columnsWidth,
			padding,
			spacing))

	// Iterate and emit rows
	for _, _row := range _t.Rows {
		// Emit row
		fmt.Println(
			formatTableLine(
				_row,
				_columnsWidth,
				padding,
				spacing))
	}
	return nil
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
func getColumnsWidth(_t table) ([]int, error) {

	_output := make([]int, len(_t.Headers))

	// Measure header
	for i, _header := range _t.Headers {
		_output[i] = getRuneCount(_header)
	}

	// Iterate and measure rows
	for _, _row := range _t.Rows {
		for j, _cell := range _row {
			if _thisLength := getRuneCount(_cell); _thisLength > _output[j] {
				_output[j] = _thisLength
			}
		}
	}
	return _output, nil
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
