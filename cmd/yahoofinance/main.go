// Command line interface to fetch quotes for equities and currencies from
// Yahoo Finance.
package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/scottjbarr/yahoofinance"
)

func usage() {
	fmt.Printf("Usage : %s symbol [symbol] ...\n", os.Args[0])
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	// get symbols from the command line
	symbols := os.Args[1:len(os.Args)]
	sort.Strings(symbols)

	client := yahoofinance.NewClient()
	quotes, err := client.GetQuotes(symbols)

	if err != nil {
		fmt.Println("%+v\n", err)
	}

	for _, quote := range quotes {
		fmt.Println(quote)
	}
}
