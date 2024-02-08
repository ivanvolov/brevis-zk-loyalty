package age

import (
	"context"
	"github.com/brevis-network/brevis-sdk/sdk"
	"github.com/brevis-network/brevis-sdk/test"
	"github.com/ethereum/go-ethereum/common"
	"testing"
)

func TestCircuit(t *testing.T) {
	app, err := sdk.NewBrevisApp("") // TODO use your eth rpc
	check(err)

	txHash := common.HexToHash(
		"8b805e46758497c6b32d0bf3cad3b3b435afeb0adb649857f24e424f75b79e46")

	app.AddTransaction(sdk.TransactionQuery{TxHash: txHash})

	guest := &AppCircuit{}
	guestAssignment := &AppCircuit{}

	circuitInput, err := app.BuildCircuitInput(context.Background(), guest)
	check(err)

	test.ProverSucceeded(t, guest, guestAssignment, circuitInput)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
