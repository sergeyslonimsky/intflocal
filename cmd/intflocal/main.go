package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/sergeyslonimsky/intflocal/analyzer"
)

func main() {
	singlechecker.Main(analyzer.New(analyzer.Settings{}))
}
