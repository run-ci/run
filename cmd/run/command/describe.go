package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/run-ci/run"
)

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use:   "describe TASK",
	Short: "Describe one or more tasks",
	Run: func(cmd *cobra.Command, args []string) {
		for _, task := range args {
			err := describeTask(task)
			if err != nil {
				fmt.Printf("error describing task '%s': %s\n", task, err.Error())
			}
		}
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(describeCmd)
}

// describeTask takes the name of a task and prints the contents of that task
// for viewing in the command-line.
func describeTask(taskName string) error {
	task, err := run.LoadTask(taskName)
	if err != nil {
		// unable to load task
		return err
	}

	header := fmt.Sprintf("Description of Task '%s'", task.Name)
	divider := strings.Repeat("=", len(header))
	bounds := strings.Repeat("=", 64)
	fmt.Printf("%s\n%s\n%s\n\n", bounds, header, divider)

	desc := task.Description
	if desc[len(desc)-1] != '\n' {
		// found a line without a newline character in it
		desc = fmt.Sprintf("%v\n", desc)
	}

	fmt.Println(desc)

	if len(task.Arguments) > 0 {
		header = fmt.Sprintf("Arguments for Task '%s'", task.Name)
		divider = strings.Repeat("=", len(header))
		fmt.Printf("%s\n%s\n\n", header, divider)

		for k, arg := range task.Arguments {
			line := k
			if arg.Description != "" {
				line = fmt.Sprintf("%v: %v", line, arg.Description)
			}

			if arg.Default != "" {
				line = fmt.Sprintf("%s (defaults to %v)", line, arg.Default)
			}

			if line[len(line)-1] != '\n' {
				// found a line without a newline character in it
				line = fmt.Sprintf("%v\n", line)
			}

			fmt.Print(line)
		}
	}

	fmt.Println(bounds)
	return nil
}
