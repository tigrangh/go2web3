package alchemy

import (
	"encoding/json"
	"errors"
	"fmt"
	"goland/go2web3/go2web3common"
)

func alchemyParam(transactionData go2web3common.TransactionData) map[string]interface{} {
	return map[string]interface{}{
		"from":     transactionData.FromAddress,
		"to":       transactionData.InteractAddress,
		"value":    fmt.Sprintf("0x%x", transactionData.EtherValue),
		"data":     fmt.Sprintf("0x%x", transactionData.CallData),
		"gas":      fmt.Sprintf("0x%x", transactionData.GasLimit),
		"gasPrice": fmt.Sprintf("0x%x", transactionData.GasPrice)}
}

func DecideAlchemyEndpoint(chainId uint) string {
	if chainId == 1 {
		return "eth-mainnet"
	} else if chainId == 137 {
		return "polygon-mainnet"
	} else if chainId == 5 {
		return "eth-goerly"
	} else if chainId == 80001 {
		return "polygon-mumbai"
	} else {
		panic(errors.New(fmt.Sprintf("alchemy endpoint undefined for %d", chainId)))
	}
}

func PrintAssetChangeRequest(chainId uint, transactionData go2web3common.TransactionData) {
	endpointPart := DecideAlchemyEndpoint(chainId)

	param := alchemyParam(transactionData)

	params := []interface{}{param}

	jsonrpc := map[string]interface{}{
		"id":      1,
		"jsonrpc": "2.0",
		"method":  "alchemy_simulateAssetChanges",
		"params":  params}

	postData, _ := json.Marshal(jsonrpc)

	curl := fmt.Sprintf(`curl https://%s.g.alchemy.com/v2/docs-demo -H 'Origin: https://docs.alchemy.com' -H 'accept: application/json' -H 'content-type: application/json' --data-binary '%s'`, endpointPart, postData)

	fmt.Println("run the following command in terminal to send request to alchemy")
	fmt.Println(curl)
}

func PrintAssetChangeBundleRequest(chainId uint, transactionsData []go2web3common.TransactionData) {
	endpointPart := DecideAlchemyEndpoint(chainId)

	param := make([]interface{}, len(transactionsData))
	for index, transactionData := range transactionsData {
		param[index] = alchemyParam(transactionData)
	}

	params := []interface{}{param}

	jsonrpc := map[string]interface{}{
		"id":      1,
		"jsonrpc": "2.0",
		"method":  "alchemy_simulateAssetChangesBundle",
		"params":  params}

	postData, _ := json.Marshal(jsonrpc)

	curl := fmt.Sprintf(`curl https://%s.g.alchemy.com/v2/docs-demo -H 'Origin: https://docs.alchemy.com' -H 'accept: application/json' -H 'content-type: application/json' --data-binary '%s'`, endpointPart, postData)

	fmt.Println("run the following command in terminal to send request to alchemy")
	fmt.Println(curl)
}
