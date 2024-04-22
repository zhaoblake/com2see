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
	var rt string

	flag.StringVar(&path, "path", "", "Your git repository path, defaults to current working directory")
	flag.IntVar(&year, "year", time.Now().Year()-1, "Year for the commit report, defaults to last year")
	flag.StringVar(&rt, "rt", "console", "Report type, defaults to console output report")
	flag.Parse()

	var err error
	switch rt {
	case "html":
		err = internal.GenerateHTMLReport(path, year)
	case "console":
		err = internal.GenerateConsoleReport(path, year)
	default:
		fmt.Println("Wrong report type!")
		os.Exit(1)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
