package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/run-ci/run"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "run [flags] -- TASK",
	Short: "run your tasks from anywhere",
	Run: func(cmd *cobra.Command, args []string) {
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
