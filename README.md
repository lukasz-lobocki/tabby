# tabby ![Static](https://img.shields.io/badge/kopiec-majowy-honeydew?style=for-the-badge&labelColor=floralwhite)

Prints left-aligned table.

ANSI color sequences within cells **do not** distort the alignment.

## Usage

```golang
package main

import 	"github.com/lukasz-lobocki/tabby"

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
		"eins \033[4;33mzwei\033[0m drei",
		"vier",
	}); err != nil {
		log.Fatalln(err)
	}

	if err := table.AppendRow(tabby.Row{
		"bądź co nieco",
		"Będzin \033[0;31mwątły królik\033[0m mąką",
	}); err != nil {
		log.Fatalln(err)
	}

	table.Print(nil)
	// alternatively, with Config provided:
	// table.Print(&tabby.Config{Spacing: "|", Padding: "."})
}
```

## Result

![Alt text](<out.png>)

## Installation

```bash
go get github.com/lukasz-lobocki/tabby@latest
```

## Credits

Inspired by [table](https://github.com/tomlazar/table) and [tabby](https://github.com/cheynewallace/tabby).
