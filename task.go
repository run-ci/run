package main

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// Task is the combination of a thing to do with metadata about
// what it does.
type Task struct {
	Name string // Name comes from the name of the file.

	Summary     string `yaml:"summary"`
	Description string `yaml:"description"`

	Image     string         `yaml:"image"`
	Command   string         `yaml:"command"`
	Mount     string         `yaml:"mount"`
	Shell     string         `yaml:"shell"` // Any shell that can take -c to execute commands.
	Arguments map[string]Arg `yaml:"arguments"`
}

// Arg is a parameter passed to the task.
type Arg struct {
	Description string `yaml:"description"`
	Default     string `yaml:"default"`
}

// LoadTask loads a task from a YAML file and returns it.
func LoadTask(name string) (Task, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return Task{}, err
	}

	printDebug("looking for tasks in %v", pwd)

	path := fmt.Sprintf("%v/tasks/%v.yaml", pwd, name)
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return Task{}, err
	}

	task := Task{Name: name}
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

// GetCmd returns the task's command as a CMD for the Docker
// container to run.
func (t Task) GetCmd() []string {
	printDebug("task %v - getting cmd", t.Name)

	// Docker won't do environment variable substitution when
	// CMD is passed as an array. This is their suggested
	// workaround in the docs.
	return []string{t.Shell, "-c", t.Command}
}

// GetEnv returns the task's arguments as ENVs for the Docker
// container to run.
func (t Task) GetEnv() ([]string, error) {
	env := []string{}

	for k, arg := range t.Arguments {
		val := arg.Default
		override := os.Getenv(k)
		if override != "" {
			val = override
		}

		if val == "" {
			return []string{}, fmt.Errorf("argument %v empty", k)
		}

		env = append(env, fmt.Sprintf("%v=%v", k, val))
	}

	return env, nil
}
