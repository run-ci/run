package run

import (
	"fmt"
	"os"
)

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

// GetCmd returns the task's command as a CMD for the Docker
// container to run.
func (t Task) GetCmd() []string {
	// Docker won't do environment variable substitution when
	// CMD is passed as an array. This is their suggested
	// workaround in the docs.
	return []string{t.Shell, "-c", t.Command}
}

// GetEnv returns the task's arguments as key-value pairs.
// If the environment specifies an argument, that value is
// used instead of the default value.
func (t Task) GetEnv() (map[string]string, error) {
	env := map[string]string{}

	for k, arg := range t.Arguments {
		val := arg.Default
		override := os.Getenv(k)
		if override != "" {
			val = override
		}

		if val == "" {
			return map[string]string{}, fmt.Errorf("argument %v empty", k)
		}

		env[k] = val
	}

	return env, nil
}
