package main

import (
	"com2see/internal"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	var path string
	var year int

	flag.StringVar(&path, "path", "", "Your git repository path, defaults to current working directory")
	flag.IntVar(&year, "year", time.Now().Year()-1, "Year for the commit report, defaults to last year")

	flag.Parse()
	err := internal.GenerateHTMLReport(path, year)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
