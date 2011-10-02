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
	for _, e := range cl {
		norm, ok := e.Data.(Normal)
		if !ok {
			log.Printf("event not normal: %#v", e)
			continue
		}
		src, dst := norm.GetSource(), norm.GetDest()
		categorize(src)
		categorize(dst)
	}
	log.Printf("Processed %d log entries with %d units in %d groups", len(cl), len(seen), len(groups))

	for group, units := range groups {
		fmt.Printf("Group 0x%04x: (%d/%d)\n", 1 << uint(group), len(units), len(seen))
		for _, unit := range units {
			fmt.Printf(" - %15b %s\n", unit.Flags, unit.Name)
		}
	}
}

type Normal interface {
	GetSource() combatlog.Unit
	GetDest() combatlog.Unit
}

var seen = map[string]bool{}
var groups = [31][]combatlog.Unit{}

func categorize(unit combatlog.Unit) {
	if _, ok := seen[unit.Name]; ok {
		return
	}
	seen[unit.Name] = true
	for shift := range groups {
		mask := combatlog.UnitFlags(1) << uint(shift)
		if unit.Flags & mask != 0 {
			groups[shift] = append(groups[shift], unit)
		}
	}
}
