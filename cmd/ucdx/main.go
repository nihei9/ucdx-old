package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "ucdx",
	Short:         "UCD (Unicode Character Database) utilities",
	SilenceErrors: true,
	SilenceUsage:  true,
}

func execute() int {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(execute())
}
