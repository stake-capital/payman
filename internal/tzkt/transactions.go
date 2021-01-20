package tzkt

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

/*
PastTransaction -
see: https://api.tzkt.io/#operation/Accounts_GetOperations
*/
type PastTransaction struct {
	Type      string    `json:"type"`
	ID        int       `json:"id"`
	Level     int       `json:"level"`
	Timestamp time.Time `json:"timestamp"`
	Block     string    `json:"block"`
	Hash      string    `json:"hash"`
	Counter   int       `json:"counter"`
	Sender    struct {
		Name    string `json:"alias"`
		Address string `json:"address"`
	} `json:"sender"`
	GasLimit      int `json:"gasLimit"`
	GasUsed       int `json:"gasUsed"`
	StorageLimit  int `json:"storageLimit"`
	StorageUsed   int `json:"storageUsed"`
	BakerFee      int `json:"bakerFee"`
	StorageFee    int `json:"storageFee"`
	AllocationFee int `json:"allocationFee"`
	Target        struct {
		Address string `json:"address"`
	} `json:"target"`
	Amount       int    `json:"amount"`
	Status       string `json:"status"`
	HasInternals bool   `json:"hasInternals"`
}

/*
GetPastTransactionsByHash -
see: https://api.tzkt.io/#operation/Operations_GetTransactionByHash
*/
func (t *Tzkt) GetPastTransactionsByHash(hash []string) ([]PastTransaction, error) {
	result := make([]PastTransaction, 0)

	for _, s := range hash {
		fmt.Printf("%s\n", s)
		resp, err := t.get(fmt.Sprintf("/v1/operations/transactions/%s", s))

		if err != nil {
			return []PastTransaction{}, errors.Wrapf(err, "failed to get transactions")
		}

		var transactions []PastTransaction
		if err := json.Unmarshal(resp, &transactions); err != nil {
			return []PastTransaction{}, errors.Wrap(err, "failed to get transactions")
		}

		result = append(result, transactions...)
	}

	return result, nil
}
