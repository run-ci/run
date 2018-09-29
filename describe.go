package main

import (
	"fmt"
	"strings"
)

func runDescribe(name string) {
	task, err := LoadTask(name)
	if err != nil {
		printDebug("got error %v for task %v, skipping...", err, name)
	}

	header := fmt.Sprintf("Description of %v", task.Name)
	divider := strings.Repeat("=", len(header))

	fmt.Println(header)
	fmt.Println(divider)
	fmt.Println()

	line := task.Description
	if line[len(line)-1] != '\n' {
		printDebug("found line without newline")
		line = fmt.Sprintf("%v\n", line)
	}

	fmt.Println(line)

	fmt.Println("## Arguments")
	for k, arg := range task.Arguments {
		line := k

		if arg.Description != "" {
			line = fmt.Sprintf("%v: %v", line, arg.Description)
		}

		if line[len(line)-1] != '\n' {
			printDebug("found line without newline")
			line = fmt.Sprintf("%v\n", line)
		}

		fmt.Print(line)
	}

	fmt.Println()
}
