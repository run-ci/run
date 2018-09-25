package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

var verbose bool

func main() {
	var list bool

	flag.BoolVarP(&list, "list", "l", false, "list all tasks")
	flag.BoolVarP(&verbose, "verbose", "v", false, "print debug output")
	flag.Parse()

	switch {
	case list:
		printDebug("list called")
		err := runList()
		if err != nil {
			printFatal("error listing tasks: %v", err)
		}
	}
}

func printDebug(msg string, args ...interface{}) {
	if verbose {
		msg = fmt.Sprintf("[DEBUG] %v\n", msg)

		fmt.Printf(msg, args...)
	}
}

func printFatal(msg string, args ...interface{}) {
	msg = fmt.Sprintf("%v\n", msg)
	fmt.Printf(msg, args)

	os.Exit(1)
}
