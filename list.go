package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func runList() error {
	files, err := ioutil.ReadDir("./tasks")
	if err != nil {
		printDebug("error listing files: %v", err)
		return err
	}

	for _, finfo := range files {
		printDebug("found file: %v", finfo.Name())

		name := strings.Split(finfo.Name(), ".")[0]

		task, err := LoadTask(name)
		if err != nil {
			printDebug("got error %v for task %v, skipping...", err, name)
		}

		line := fmt.Sprintf("%v: %v", task.Name, task.Summary)
		if line[len(line)-1] != '\n' {
			printDebug("found line without newline")
			line = fmt.Sprintf("%v\n", line)
		}

		fmt.Print(line)
	}

	return nil
}
