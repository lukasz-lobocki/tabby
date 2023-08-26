package main

import (
	"log"

	"github.com/lukasz-lobocki/tabby/pkg/tabby"
	"github.com/lukasz-lobocki/tabby/pkg/utils"
)

func main() {

	_table := new(tabby.Table)

	if err := _table.SetHeader(
		tabby.Header{
			"Something",
			"Another",
			"One mo" + utils.RED + "r" + utils.RESET + "e",
		},
	); err != nil {
		log.Fatalln(err)
	}

	if err := _table.AppendRow(
		tabby.Row{
			"uno",
			"dos",
			"tres",
			//"quatrro",
		},
	); err != nil {
		log.Fatalln(err)
	}

	if err := _table.AppendRow(
		tabby.Row{
			"first",
			"seco" + utils.UNDERLINE_GREEN + "n" + utils.RESET + "d_garbage67890",
			"third",
		},
	); err != nil {
		log.Fatalln(err)
	}

	_table.Print(nil)
}
