package main

import (
	"fmt"
	"os"

	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/sergeyslonimsky/intflocal/analyzer"
)

var version = "dev"

func main() {
	if len(os.Args) == 2 && (os.Args[1] == "-version" || os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Println("intflocal version " + version)
		os.Exit(0)
	}

	singlechecker.Main(analyzer.New(analyzer.Settings{}))
}
