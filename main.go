package main

import (
	"fmt"
	"os"

	"github.com/srackham/xrate/internal/xrates"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Print(`
A simple CLI command to print the amount of CURRENCY that $1 USD would buy at today's rates.

Usage: xrate CURRENCY

CURRENCY is the fiat currency's ticker symbol e.g. NZD, AUD, EUR
`)
		os.Exit(1)
	}

	currency := os.Args[1]

	x := xrates.New()
	err := x.Load()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
	defer x.Save()
	rate, err := x.GetCachedRate(currency, false)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Printf("%.2f\n", rate)
}
