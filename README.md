# go-hyperliquid
 A golang SDK for [Hyperliquid PerpDEX](https://hyperliquid.xyz/).

# API reference
- [Hyperliquid API docs](https://hyperliquid.gitbook.io/hyperliquid-docs/for-developers/api)
- [Hyperliquid official Python SDK](https://github.com/hyperliquid-dex/hyperliquid-python-sdk)

# How to install?
```
go get github.com/Logarithm-Labs/go-hyperliquid/hyperliquid
```

### Documentation

[![GoDoc](https://godoc.org/github.com/adshao/go-binance?status.svg)](https://pkg.go.dev/github.com/Logarithm-Labs/go-hyperliquid/hyperliquid#section-documentation)


# Quick start
```
package main

import (
	"log"

	"github.com/Logarithm-Labs/go-hyperliquid/hyperliquid"
)

func main() {
	hyperliquidClient := hyperliquid.NewHyperliquid(&hyperliquid.HyperliquidClientConfig{
		IsMainnet:      true,
		AccountAddress: "0x12345",   // Main address of the Hyperliquid account that you want to use
		PrivateKey:     "abc1234",   // Private key of the account or API private key from Hyperliquid
	})

	// Get balances
	res, err := hyperliquidClient.GetAccountState()
	if err != nil {
		log.Print(err)
	}
	log.Printf("GetAccountState(): %+v", res)
}
```
