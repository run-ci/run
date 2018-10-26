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

var (
	flagDescribeTask *bool
	flagListTasks    *bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "run [flags] TASK",
	Short: "run your tasks from anywhere",
	Run: func(cmd *cobra.Command, args []string) {
		if *flagListTasks {
			listTasks(taskDir)
			os.Exit(0)
		}

		if *flagDescribeTask {
			if len(args) == 0 {
				fmt.Println("must provide one or more tasks to describe")
				fmt.Println(cmd.Usage())
				os.Exit(1)
			}
			for _, task := range args {
				err := describeTask(task)
				if err != nil {
					fmt.Printf("error describing task '%s': %s\n", task, err)
				}
			}
			os.Exit(0)
		}

		if len(args) == 0 {
			fmt.Println("must provide one or more tasks to run")
			fmt.Println(cmd.Usage())
			os.Exit(1)
		}

		for _, task := range args {
			run.RunTask(task)
		}
		os.Exit(0)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	flagListTasks = rootCmd.Flags().BoolP("list", "l", false, "list all tasks")
	flagDescribeTask = rootCmd.Flags().BoolP("describe", "d", false, "describe one or more tasks")
}

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

		line := fmt.Sprintf("%v: %v", task.Name, task.Summary)
		if line[len(line)-1] != '\n' {
			// found a line without a newline character in it
			line = fmt.Sprintf("%v\n", line)
		}

		fmt.Print(line)
	}

	return nil
}

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
