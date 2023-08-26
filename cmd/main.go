package main

import (
	"log"

	"github.com/lukasz-lobocki/tabby/pkg/tabby"
	"github.com/lukasz-lobocki/tabby/pkg/utils"
)

func main() {

	_tab := new(tabby.Table)

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
