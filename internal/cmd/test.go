package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// TestCommand returns the cobra command for test
func TestCommand() *cobra.Command {
	var test = &cobra.Command{
		Use:     "test",
		Short:   "test just shows random text",
		Long:    "test just shows random text",
		Example: `tzpay test`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("Testing cron job")
		},
	}

	return test
}
