package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	g_input = flag.String("f", "", "file to parse")
)

func show_usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s -f file<file>\n\n",
		os.Args[0])
	fmt.Fprintf(os.Stderr,
		"Flags:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = show_usage
	flag.Parse()
	var retval int
	if len(*g_input) == 0 {
		retval = 1
	} else {

		retval = parseFile(*g_input)
	}
	if retval != 0 {
		flag.Usage()
	}
	os.Exit(retval)
}
