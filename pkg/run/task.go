package run

import (
	"fmt"
	"os"
	"strings"

	vault "github.com/hashicorp/vault/api"
)

type Task struct {
	Name string `yaml:"name" json:"name"`

	Summary     string `yaml:"summary" json:"summary"`
	Description string `yaml:"description" json:"description"`

	Image     string         `yaml:"image" json:"image"`
	Command   string         `yaml:"command" json:"command"`
	Mount     string         `yaml:"mount" json:"mount"`
	Shell     string         `yaml:"shell" json:"shell"` // Any shell that can take -c to execute commands.
	Arguments map[string]Arg `yaml:"arguments" json:"arguments"`
}

// Arg is a parameter passed to the task.
type Arg struct {
	Description string `yaml:"description" json:"description"`
	Default     string `yaml:"default" json:"default"`
	Vault       string `yaml:"vault" json:"vault"`
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
// If the argument is specified as a vault path (using
// Arg.Vault), the value at that path is used instead of
// the default value.
// If the environment specifies an argument, that value is
// used instead of the default value or the Vault value.
func (t Task) GetEnv() (map[string]string, error) {
	env := map[string]string{}

	for k, arg := range t.Arguments {
		val := arg.Default

		if arg.Vault != "" {
			vaultcfg := vault.DefaultConfig()
			if vaultcfg.Error != nil {
				return nil, vaultcfg.Error
			}

			client, err := vault.NewClient(vaultcfg)
			if err != nil {
				return nil, err
			}

			segs := strings.Split(arg.Vault, ":")

			secret, err := client.Logical().Read(segs[0])
			if err != nil {
				return nil, err
			}

			val = secret.Data[segs[1]].(string)
		}

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
