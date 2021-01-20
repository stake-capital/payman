package main

import (
	"github.com/goat-systems/tzpay/v3/internal/cmd"
	"github.com/spf13/cobra"
)

func main() {

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
		cmd.DryRun2Command(), // dryrun2
		cmd.Run2Command(),    // run2
	)

	rootCommand.Execute()
}

/**
Comment

Dryrun2 - displays two table
	1. Full Delegators list with their rewards for a given cycle	-- the same result of <dryrun>
	2. Unpaid delegators -- fetch all transactions using the hashes saved in <data/past_cycle_hash.json>

Run2 - handle payment for dryrun2
*/
