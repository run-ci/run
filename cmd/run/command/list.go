package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/run-ci/run"
)

const (
	taskDir = "tasks/"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available tasks",
	Run: func(cmd *cobra.Command, args []string) {
		listTasks(taskDir)
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

// listTasks takes the name of a directory in which your task yml files are
// and iterates over each one, printing the name and summary of each task.
func listTasks(taskDir string) error {
	taskFiles, err := ioutil.ReadDir(taskDir)
	if err != nil {
		return err
	}

	for _, fileInfo := range taskFiles {
		log.Debug("found file: %s", fileInfo.Name())
		taskName := strings.Split(fileInfo.Name(), ".")[0]

		task, err := run.LoadTask(taskName)
		if err != nil {
			log.Errorf("received an error for task %s. skipping... (err: %v)\n", taskName, err)
		}

		line := fmt.Sprintf("%12s\t%v", task.Name, task.Summary)
		if line[len(line)-1] != '\n' {
			// found a line without a newline character in it
			line = fmt.Sprintf("%v\n", line)
		}

		fmt.Print(line)
	}

	return nil
}
