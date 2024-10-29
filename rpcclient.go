// rpcclient.go
package main

import (
	"fmt"
	"log"

	"github.com/btcsuite/btcd/rpcclient"
)

var rpcClient *rpcclient.Client

func initRPCClient() {
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	connCfg := &rpcclient.ConnConfig{
		Host:         fmt.Sprintf("%s:%s", config.RPCHost, config.RPCPort),
		User:         config.RPCUser,
		Pass:         config.RPCPassword,
		HTTPPostMode: true, // Bitcoin Core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin Core does not provide TLS by default
	}

	// For HTTPS connections (if TLS is enabled on your node)
	// if !connCfg.DisableTLS {
	// 	connCfg.TLSConfig = &tls.Config{
	// 		InsecureSkipVerify: true, // Use proper certificate validation in production
	// 	}
	// }

	rpcClient, err = rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatalf("Error creating RPC client: %v", err)
	}

	// Ensure the client is closed on application exit
	// Defer the shutdown to the end of main
}
