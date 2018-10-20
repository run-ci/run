package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/run-ci/run/pkg/run"
	yaml "gopkg.in/yaml.v2"
)

// LoadTask loads a task from a YAML file and returns it.
func LoadTask(name string) (run.Task, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return run.Task{}, err
	}

	printDebug("looking for tasks in %v", pwd)

	path := fmt.Sprintf("%v/tasks/%v.yaml", pwd, name)
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return run.Task{}, err
	}

	task := run.Task{Name: name}
	err = yaml.UnmarshalStrict(f, &task)
	if err != nil {
		return task, err
	}

	if task.Mount == "" {
		task.Mount = "/mnt/repo"
	}

	if task.Shell == "" {
		task.Shell = "sh"
	}

	return task, nil
}
