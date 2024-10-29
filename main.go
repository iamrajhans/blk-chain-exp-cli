package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/go-resty/resty/v2"
	"github.com/olekukonko/tablewriter"
)

func main() {
	// Define subcommands
	blockCmd := flag.NewFlagSet("block", flag.ExitOnError)
	txCmd := flag.NewFlagSet("tx", flag.ExitOnError)
	addressCmd := flag.NewFlagSet("address", flag.ExitOnError)
	statsCmd := flag.NewFlagSet("stats", flag.ExitOnError)

	// Define flags for 'block' command
	blockHash := blockCmd.String("hash", "", "Block hash")
	blockHeight := blockCmd.Int64("height", -1, "Block height")

	// Define flags for 'tx' command
	txID := txCmd.String("txid", "", "Transaction ID (hash)")

	// Define flags for 'address' command
	address := addressCmd.String("address", "", "Bitcoin address")

	// No flags for 'stats' command at this point

	if len(os.Args) < 2 {
		fmt.Println("Expected 'block', 'tx', 'address' or 'stats' subcommands")
		os.Exit(1)
	}

	// Switch on the subcommand
	switch os.Args[1] {
	case "block":
		blockCmd.Parse(os.Args[2:])
		getBlockInfo(*blockHash, *blockHeight)
	case "tx":
		txCmd.Parse(os.Args[2:])
		getTransactionInfo(*txID)
	case "address":
		addressCmd.Parse(os.Args[2:])
		getAddressInfo(*address)
	case "stats":
		statsCmd.Parse(os.Args[2:])
		getNetworkStats()
	default:
		fmt.Println("Expected 'block', 'tx', 'address' or 'stats' subcommands")
		os.Exit(1)
	}

	// Initialize the RPC client
	// initRPCClient()
	// defer rpcClient.Shutdown()
	// initCache()
}

func getBlockInfo(hash string, height int64) {
	if hash == "" && height == -1 {
		fmt.Println("Please provide a block hash or height using --hash or --height.")
		return
	}

	var blockHash *chainhash.Hash
	var err error

	if hash != "" {
		blockHash, err = chainhash.NewHashFromStr(hash)
		if err != nil {
			fmt.Printf("Invalid block hash: %v\n", err)
			return
		}
	} else {
		blockHash, err = rpcClient.GetBlockHash(height)
		if err != nil {
			fmt.Printf("Error getting block hash by height: %v\n", err)
			return
		}
	}

	blockVerbose, err := rpcClient.GetBlockVerbose(blockHash)
	if err != nil {
		fmt.Printf("Error getting block: %v\n", err)
		return
	}

	displayBlockInfo(blockVerbose)
}

func displayBlockInfo(block *btcjson.GetBlockVerboseResult) {
	data := [][]string{
		{"Hash", block.Hash},
		{"Confirmations", fmt.Sprintf("%d", block.Confirmations)},
		{"Size", fmt.Sprintf("%d bytes", block.Size)},
		{"Height", fmt.Sprintf("%d", block.Height)},
		{"Version", fmt.Sprintf("%d", block.Version)},
		{"Merkle Root", block.MerkleRoot},
		{"Time", fmt.Sprintf("%d", block.Time)},
		{"Nonce", fmt.Sprintf("%d", block.Nonce)},
		{"Difficulty", fmt.Sprintf("%f", block.Difficulty)},
		{"Previous Hash", block.PreviousHash},
		{"Next Hash", block.NextHash},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field", "Value"})

	for _, v := range data {
		table.Append(v)
	}

	table.Render()

	fmt.Println("\nTransactions:")
	for _, tx := range block.Tx {
		fmt.Printf("  %s\n", tx)
	}
}

func getTransactionInfo(txidStr string) {
	if txidStr == "" {
		fmt.Println("Please provide a transaction ID using --txid.")
		return
	}

	txHash, err := chainhash.NewHashFromStr(txidStr)
	if err != nil {
		fmt.Printf("Invalid transaction ID: %v\n", err)
		return
	}

	rawTx, err := rpcClient.GetRawTransactionVerbose(txHash)
	if err != nil {
		fmt.Printf("Error getting transaction: %v\n", err)
		return
	}

	displayTransactionInfo(rawTx)
}

func displayTransactionInfo(tx *btcjson.TxRawResult) {
	fmt.Printf("Transaction ID: %s\n", tx.Txid)
	fmt.Printf("Hash: %s\n", tx.Hash)
	fmt.Printf("Size: %d bytes\n", tx.Size)
	fmt.Printf("Version: %d\n", tx.Version)
	fmt.Printf("LockTime: %d\n", tx.LockTime)
	fmt.Printf("Confirmations: %d\n", tx.Confirmations)
	fmt.Printf("Block Hash: %s\n", tx.BlockHash)
	fmt.Printf("Time: %d\n", tx.Time)
	fmt.Printf("Inputs:\n")
	for _, vin := range tx.Vin {
		fmt.Printf("  TXID: %s, Vout: %d\n", vin.Txid, vin.Vout)
	}
	fmt.Printf("Outputs:\n")
	for _, vout := range tx.Vout {
		fmt.Printf("  Value: %f BTC, Addresses: %v\n", vout.Value, vout.ScriptPubKey.Addresses)
	}
}

func getAddressInfo(address string) {
	if address == "" {
		fmt.Println("Please provide an address using --address.")
		return
	}

	client := resty.New()
	apiURL := fmt.Sprintf("https://api.blockcypher.com/v1/btc/main/addrs/%s/full", address)

	resp, err := client.R().Get(apiURL)
	if err != nil {
		fmt.Printf("Error fetching address info: %v\n", err)
		return
	}

	if resp.StatusCode() != 200 {
		fmt.Printf("API request failed with status: %s\n", resp.Status())
		return
	}

	fmt.Println(string(resp.Body()))
}

func getNetworkStats() {
	networkInfo, err := rpcClient.GetNetworkInfo()
	if err != nil {
		fmt.Printf("Error getting network info: %v\n", err)
		return
	}

	blockchainInfo, err := rpcClient.GetBlockChainInfo()
	if err != nil {
		fmt.Printf("Error getting blockchain info: %v\n", err)
		return
	}

	fmt.Printf("Network Version: %d\n", networkInfo.Version)
	fmt.Printf("Protocol Version: %d\n", networkInfo.ProtocolVersion)
	fmt.Printf("Connections: %d\n", networkInfo.Connections)
	fmt.Printf("Difficulty: %f\n", blockchainInfo.Difficulty)
	fmt.Printf("Chain: %s\n", blockchainInfo.Chain)
	fmt.Printf("Blocks: %d\n", blockchainInfo.Blocks)
	fmt.Printf("Best Block Hash: %s\n", blockchainInfo.BestBlockHash)
	fmt.Printf("Median Time: %d\n", blockchainInfo.MedianTime)
}
