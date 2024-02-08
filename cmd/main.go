package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/brevis-network/brevis-quickstart/age"
	"github.com/brevis-network/brevis-sdk/sdk"
	"github.com/ethereum/go-ethereum/common"
	"path/filepath"
)

var mode = flag.String("mode", "", "compile or prove")
var outDir = flag.String("out", "$HOME/circuitOut/myBrevisApp", "compilation output dir")
var srsDir = flag.String("srs", "$HOME/kzgsrs", "where to cache kzg srs")
var tx = flag.String("tx", "", "tx hash to prove")

func main() {
	flag.Parse()
	switch *mode {
	case "compile":
		compile()
	case "prove":
		prove()
	default:
		panic(fmt.Errorf("unsupported mode %s", *mode))
	}
}

func compile() {
	// This first part is copied from app/circuit_test.go. We added the source data, then we generated the circuit input.
	app, err := sdk.NewBrevisApp("https://eth-mainnet.nodereal.io/v1/0af795b55d124a61b86836461ece1dee") // TODO use your eth rpc
	check(err)
	txHash := common.HexToHash(
		"8b805e46758497c6b32d0bf3cad3b3b435afeb0adb649857f24e424f75b79e46")
	app.AddTransaction(sdk.TransactionQuery{TxHash: txHash})
	appCircuit := &age.AppCircuit{}

	circuitInput, err := app.BuildCircuitInput(context.Background(), appCircuit)
	check(err)

	// The compilation output is the description of the circuit's constraint system.
	// You should use sdk.WriteTo to serialize and save your circuit so that it can
	// be used in the proving step later.
	compiledCircuit, err := sdk.Compile(appCircuit, circuitInput)
	check(err)
	err = sdk.WriteTo(compiledCircuit, filepath.Join(*outDir, "compiledCircuit"))
	check(err)

	// Setup is a one-time effort per circuit. A cache dir can be provided to output
	// external dependencies. Once you have the verifying key you should also save
	// its hash in your contract so that when a proof via Brevis is submitted
	// on-chain you can verify that Brevis indeed used your verifying key to verify
	// your circuit computations
	pk, vk, err := sdk.Setup(compiledCircuit, *srsDir)
	check(err)
	err = sdk.WriteTo(pk, filepath.Join(*outDir, "pk"))
	check(err)
	err = sdk.WriteTo(vk, filepath.Join(*outDir, "vk"))
	check(err)

	fmt.Println("compilation/setup complete")
}

func prove() {
	if len(*tx) == 0 {
		panic("-tx is required")
	}

	// Loading the previous compile result into memory
	fmt.Println(">> Reading circuit, pk, and vk from disk")
	compiledCircuit, err := sdk.ReadCircuitFrom(filepath.Join(*outDir, "compiledCircuit"))
	check(err)
	pk, err := sdk.ReadPkFrom(filepath.Join(*outDir, "pk"))
	check(err)
	vk, err := sdk.ReadVkFrom(filepath.Join(*outDir, "vk"))
	check(err)

	// Query the user specified tx
	app, err := sdk.NewBrevisApp("https://eth-mainnet.nodereal.io/v1/0af795b55d124a61b86836461ece1dee") // TODO use your eth rpc
	check(err)
	app.AddTransaction(sdk.TransactionQuery{TxHash: common.HexToHash(*tx)})

	appCircuit := &age.AppCircuit{}
	appCircuitAssignment := &age.AppCircuit{}

	// Prove
	fmt.Println(">> Proving the transaction using my circuit")
	circuitInput, err := app.BuildCircuitInput(context.Background(), appCircuit)
	check(err)
	witness, publicWitness, err := sdk.NewFullWitness(appCircuitAssignment, circuitInput)
	check(err)
	proof, err := sdk.Prove(compiledCircuit, pk, witness)
	check(err)
	err = sdk.WriteTo(proof, filepath.Join(*outDir, "proof-"+*tx))
	check(err)

	// Test verifying the proof we just generated
	err = sdk.Verify(vk, publicWitness, proof)
	check(err)

	fmt.Println(">> Initiating Brevis request")
	appContract := common.HexToAddress("0x73090023b8D731c4e87B3Ce9Ac4A9F4837b4C1bd")
	refundee := common.HexToAddress("0x164Ef8f77e1C88Fb2C724D3755488bE4a3ba4342")

	calldata, requestId, feeValue, err := app.PrepareRequest(vk, 1, 11155111, refundee, appContract)
	check(err)
	fmt.Printf("calldata %x\n", calldata)
	fmt.Printf("feeValue %d\n", feeValue)
	fmt.Printf("requestId %s\n", requestId)

	// Submit proof to Brevis
	fmt.Println(">> Submitting my proof to Brevis")
	err = app.SubmitProof(proof)
	check(err)

	// Poll Brevis gateway for query status till the final proof is submitted
	// on-chain by Brevis and your contract is called
	fmt.Println(">> Waiting for final proof generation and submission")
	tx, err := app.WaitFinalProofSubmitted(context.Background())
	check(err)
	fmt.Printf(">> Final proof submitted: tx hash %s\n", tx)

	// [Don't forget to make the transaction that pays the fee by calling Brevis.sendRequest]
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
