package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// Task is the combination of a thing to do with metadata about
// what it does.
type Task struct {
	Name string // Name comes from the name of the file.

	Summary     string `yaml:"summary"`
	Description string `yaml:"description"`

	Image   string `yaml:"image"`
	Command string `yaml:"command"`
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

// GetCmd returns the tasks command as a CMD for the Docker
// container to run.
func (t Task) GetCmd() []string {
	printDebug("task %v - getting cmd", t.Name)
	return strings.Split(t.Command, " ")
}
