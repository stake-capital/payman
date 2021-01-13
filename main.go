package main

import (
	"log"

	"github.com/goat-systems/tzpay/v3/internal/cmd"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	rootCommand := &cobra.Command{
		Use:   "tzpay",
		Short: "A bulk payout tool for bakers in the Tezos Ecosystem",
	}
	rootCommand.AddCommand(
		cmd.DryRunCommand(),
		cmd.ServCommand(),
		cmd.RunCommand(),
		cmd.NewVersionCommand(),
		cmd.NewSetupCommand(),
		cmd.DryRunOldCyclesCommand(), // dryrun2
	)

	rootCommand.Execute()
}
