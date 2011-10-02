package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kylelemons/wowlog/combatlog"
)

func main() {
	flag.Usage = func() {
		fmt.Printf(""+
`Usage:
	%s [options] <combatlog>

Options:
`, os.Args[0])
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		return
	}

	filename := args[0]

	log.Printf("Parsing %s...", filename)
	cl, err := combatlog.ReadFile(filename)
	if err != nil {
		log.Fatal("graphlog: %s", err)
	}

	log.Printf("Analyzing %d records...", len(cl))
}
