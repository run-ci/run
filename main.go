package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

var verbose bool

func main() {
	var list bool
	var taskname string

	flag.BoolVarP(&list, "list", "l", false, "list all tasks")
	flag.BoolVarP(&verbose, "verbose", "v", false, "print debug output")
	flag.StringVarP(&taskname, "describe", "d", "", "describe target task")
	flag.Parse()

	switch {
	case list:
		printDebug("list called")
		err := runList()
		if err != nil {
			printFatal("error listing tasks: %v", err)
		}

	case taskname != "":
		printDebug("describe called with taskname %v", taskname)
		runDescribe(taskname)

	default:
		taskname := os.Args[1]
		printDebug("running task %v", taskname)
		runTask(taskname)
	}
}

func printDebug(msg string, args ...interface{}) {
	if verbose {
		msg = fmt.Sprintf("%v\n", msg)

		fmt.Printf(msg, args...)
	}
}

func printFatal(msg string, args ...interface{}) {
	msg = fmt.Sprintf("%v\n", msg)
	fmt.Printf(msg, args)

	os.Exit(1)
}
