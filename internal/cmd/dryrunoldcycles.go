package cmd

import (
	"strconv"

	"github.com/goat-systems/tzpay/v3/internal/config"
	"github.com/goat-systems/tzpay/v3/internal/payout"
	"github.com/goat-systems/tzpay/v3/internal/print"
	"github.com/goat-systems/tzpay/v3/internal/tzkt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// DryRunOldCycles -
type DryRunOldCycles struct {
	payout payout.IFace
	config config.Config
	tzkt   tzkt.IFace
	cycle  int
}

// NewDryRunOldCycles returns a new DryRunOldCycles
func NewDryRunOldCycles(cycle string) DryRunOldCycles {
	config, err := config.New()
	if err != nil {
		log.WithField("error", err.Error()).Fatal("Failed to load config.")
	}

	// Clear sensitive data if loaded
	config.Key.Password = ""
	config.Key.Esk = ""

	c, err := strconv.Atoi(cycle)
	if err != nil {
		log.WithField("error", err.Error()).Fatal("Failed to parse cycle argument into integer.")
	}

	payout, err := payout.New(config, c, false, false)
	if err != nil {
		log.WithField("error", err.Error()).Fatal("Failed to intialize payout.")
	}

	return DryRunOldCycles{
		payout: payout,
		config: config,
		tzkt:   tzkt.NewTZKT(config.API.TZKT),
		cycle:  c,
	}
}

// DryRunOldCyclesCommand returns the cobra command for dryrun2
// Only table format is supported
// It should return the reward table(same as dryrun) + real transaction happened for that cycle (6 cycles after)
// The script should find delegators who are not get paid for their rewards and print those who need to be paid for past cycles
func DryRunOldCyclesCommand() *cobra.Command {

	var dryrun = &cobra.Command{
		Use:     "dryrun2",
		Short:   "dryrun2 simulates a payout",
		Long:    "dryrun2 simulates a payout and prints the result in json or a table",
		Example: `tzpay dryrun2 <cycle>`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Missing cycle as argument.")
			}

			dryrun := NewDryRunOldCycles(args[0])
			dryrun.execute()
		},
	}

	return dryrun
}

func (d *DryRunOldCycles) execute() {
	rewardsSplit, err := d.payout.Execute()
	pastTransactions, error := d.tzkt.GetPastTransactions(d.cycle, d.config.Baker.PayoutAddress)

	if err != nil {
		log.WithField("error", err.Error()).Fatal("Failed to execute payout - get Reward splits.")
	}
	if error != nil {
		log.WithField("error", error.Error()).Fatal("Failed to execute payout - past Transactions.")
	}

	print.Table(d.cycle, d.config.Baker.Address, rewardsSplit)
	print.TablePastTransactions(pastTransactions)
}
