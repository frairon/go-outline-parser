package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	input   = flag.String("f", "", "file to parse")
	verbose = flag.Bool("v", false, "verbose mode")
)

func debug(s ...interface{}) {
	if *verbose {
		fmt.Println(s...)
	}
}

func showUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s -f file<file>\n\n",
		os.Args[0])
	fmt.Fprintf(os.Stderr,
		"Flags:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = showUsage
	flag.Parse()

	if len(*input) == 0 {
		flag.Usage()
		os.Exit(1)
	} else {
		output, err := parseFile(*input)
		if err != nil {
			os.Stdout.WriteString(fmt.Sprintf(`{error: "%v"}`, err))
			os.Exit(1)
		}
		os.Stdout.WriteString(output)
		os.Exit(0)
	}
}
