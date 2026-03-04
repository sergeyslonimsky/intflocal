package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/sergeyslonimsky/intflocal/analyzer"
)

var version = "dev"

func main() {
	if len(os.Args) == 2 && (os.Args[1] == "-version" || os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Println("intflocal version " + getVersion())
		os.Exit(0)
	}

	singlechecker.Main(analyzer.New(analyzer.Settings{}))
}

func getVersion() string {
	if version != "dev" {
		return version
	}

	info, ok := debug.ReadBuildInfo()
	if ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}

	return version
}
