package run

import (
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"

	"github.com/run-ci/run/pkg/run"
)

// LoadTask is a helper function that reads a task from the
// filesystem and stores it in a Task struct.
func LoadTask(name string) (run.Task, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return run.Task{}, err
	}

	log.Debugf("looking for tasks in %v", pwd)

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
