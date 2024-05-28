package main

import "fmt"

func main() {
	sdk := newBCFWalletSDK()

	wallet := sdk.newBCFWallet("34.84.178.63", 19503, "https://qapmapi.pmchainbox.com/browser")
	balance := wallet.getAddressBalance("cEAXDkaEJgWKMM61KYz2dYU1RfuxbB8Ma", "XXVXQ", "PMC")
	fmt.Println("balance=", balance)
	defer sdk.Close()
}
