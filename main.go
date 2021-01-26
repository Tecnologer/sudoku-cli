package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tecnologer/sudoku/clients/cli/cmd"
)

var (
	verbouse   = flag.Bool("v", false, "verbouse")
	version    string
	minversion string
)

func main() {
	flag.Parse()

	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("%s%s\n", version, minversion)
		return
	}

	for {
		cmd.CallCmd("")
	}
}

func printCmds() {
	fmt.Println("Sudoku-CLI provides the following commands to start")
	fmt.Println()
	fmt.Println("  version", "\n\treturns the current version")
}
