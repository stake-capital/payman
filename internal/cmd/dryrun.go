package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	gotezos "github.com/goat-systems/go-tezos/v2"
	"github.com/goat-systems/tzpay/v3/internal/config"
	"github.com/goat-systems/tzpay/v3/internal/payout"
	"github.com/goat-systems/tzpay/v3/internal/print"
	"github.com/goat-systems/tzpay/v3/internal/tzkt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// DryRun -
type DryRun struct {
	payout payout.IFace
	config config.Config
	tzkt   tzkt.IFace
	cycle  int
	table  bool
}

// NewDryRun returns a new dryrun
func NewDryRun(cycle string, table bool) DryRun {
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

	return DryRun{
		payout: payout,
		config: config,
		cycle:  c,
		table:  table,
		tzkt:   tzkt.NewTZKT(config.API.TZKT),
	}
}

// DryRunCommand returns the cobra command for dryrun
func DryRunCommand() *cobra.Command {
	var table bool

	var dryrun = &cobra.Command{
		Use:     "dryrun",
		Short:   "dryrun simulates a payout",
		Long:    "dryrun simulates a payout and prints the result in json or a table",
		Example: `tzpay dryrun <cycle>`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Missing cycle as argument.")
			}

			dryrun := NewDryRun(args[0], table)
			dryrun.execute()
		},
	}
	dryrun.PersistentFlags().BoolVarP(&table, "table", "t", false, "formats result into a table (Default: json)")

	return dryrun
}

func (d *DryRun) execute() {
	var pastTransactions []tzkt.PastTransaction

	hashValue := getHashArrayFromCycle(d.cycle)
	if len(hashValue) != 0 {
		data, error := d.tzkt.GetPastTransactionsByHash(hashValue)
		if error != nil {
			log.WithField("error", error.Error()).Fatal("Failed to execute payout - past Transactions.")
		}
		pastTransactions = data
	}

	rewardsSplit, err := d.payout.Execute(pastTransactions)
	if err != nil {
		log.WithField("error", err.Error()).Fatal("Failed to execute payout.")
	}

	amount, whiteListed, blackListed := getPaymentInfo(rewardsSplit)
	currentBalance := d.tzkt.GetCurrentBalance(d.config.Baker.PayoutAddress)
	fmt.Printf("You have to pay %f XTZ.\n", amount)
	fmt.Printf("Current balance of the wallet: %f XTZ.\n", float64(currentBalance)/float64(gotezos.MUTEZ))
	fmt.Printf("whiteListed Addresses: %d\n", whiteListed)
	fmt.Printf("Blacklisted Addresses: %d\n", blackListed)

	if d.table {
		print.Table(d.cycle, d.config.Baker.Address, rewardsSplit)
	} else {
		err := print.JSON(rewardsSplit)
		if err != nil {
			log.WithField("error", err.Error()).Fatal("Failed to print JSON report.")
		}
	}
}

func getHashArrayFromCycle(cycle int) []string {
	file, err := ioutil.ReadFile("data/past_cycle_hash.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed read file: %s\n", err)
		os.Exit(1)
	}

	var f map[string][]string
	err = json.Unmarshal(file, &f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse JSON: %s\n", err)
		os.Exit(1)
	}

	data := f[fmt.Sprint((cycle))]
	if data == nil {
		return nil
	}
	return data
}

func getPaymentInfo(rewardSplit tzkt.RewardsSplit) (float64, int, int) {
	amount := 0.00
	whiteListed := 0
	blackListed := 0
	for _, delegator := range rewardSplit.Delegators {
		if delegator.BlackListed == false {
			amount += float64(delegator.NetRewards) / float64(gotezos.MUTEZ)
			whiteListed++
		} else {
			blackListed++
		}
	}
	return amount, whiteListed, blackListed
}
