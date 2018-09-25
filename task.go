package main

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Task is the combination of a thing to do with metadata about
// what it does.
type Task struct {
	Name string // Name comes from the name of the file.

	Summary string `yaml:"summary"`
}

// LoadTask loads a task from a YAML file and returns it.
func LoadTask(name string) (Task, error) {
	path := fmt.Sprintf("./tasks/%v.yaml", name)
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return Task{}, err
	}

	task := Task{Name: name}
	return task, yaml.UnmarshalStrict(f, &task)
}
