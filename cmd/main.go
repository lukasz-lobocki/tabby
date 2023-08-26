package main

import (
	"log"

	"github.com/lukasz-lobocki/tabby/pkg/tabby"
	"github.com/lukasz-lobocki/tabby/pkg/utils"
)

func main() {

	_tab := new(tabby.Table)

	if err := _tab.AddHeaders([]string{
		"Something",
		"Another",
		"One mo" + utils.RED + "r" + utils.RESET + "e",
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
		"first",
		"seco" + utils.RED + "n" + utils.RESET + "d_garbage67890",
		"third",
	}); err != nil {
		log.Fatalln(err)
	}

	_tab.Print(nil)
}
