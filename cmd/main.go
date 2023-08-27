package main

import (
	"log"

	"github.com/lukasz-lobocki/tabby/pkg/tabby"
)

func main() {

	table := new(tabby.Table)

	if err := table.SetHeader(tabby.Header{
		"\033[4mFIRST\033[0m",
		"\033[4mSECOND\033[0m",
	}); err != nil {
		log.Fatalln(err)
	}

	if err := table.AppendRow(tabby.Row{
		"uno",
		"dos",
	}); err != nil {
		log.Fatalln(err)
	}

	if err := table.AppendRow(tabby.Row{
		"ein \033[4;33mzwei\033[0m drei",
		"vier",
	}); err != nil {
		log.Fatalln(err)
	}

	if err := table.AppendRow(tabby.Row{
		"bądź co nieco",
		"Będzin \033[0;31mkróluje\033[0m nad Polską",
	}); err != nil {
		log.Fatalln(err)
	}

	table.Print(nil)
}
