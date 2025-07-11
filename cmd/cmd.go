package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

func Execute() error {
	rootCmd := newRootCommand()

	rootCmd.AddCommand(serverStartCommand())

	return rootCmd.Execute()
}

func newRootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "app-test",
		Short: "Test task for Golang Developer",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Usage(); err != nil {
				log.Printf("%+v\n", err)
			}
		},
	}
}
