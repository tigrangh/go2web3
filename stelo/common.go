package stelo

import (
	"encoding/json"
	"fmt"
	"goland/go2web3/go2web3common"
)

func steloParam(chainId uint, transactionData go2web3common.TransactionData) map[string]interface{} {
	return map[string]interface{}{
		"chainId":  chainId,
		"from":     transactionData.FromAddress,
		"to":       transactionData.InteractAddress,
		"data":     fmt.Sprintf("0x%x", transactionData.CallData),
		"value":    fmt.Sprintf("0x%x", transactionData.EtherValue),
		"gasLimit": fmt.Sprintf("0x%x", transactionData.GasLimit),
		"gas":      fmt.Sprintf("0x%x", transactionData.GasLimit),
		"gasPrice": fmt.Sprintf("0x%x", transactionData.GasPrice)} //"url":      fmt.Sprintf("0x%x",transactionData.)

}

func PrintAssetChangeRequest(chainId uint, transactionData go2web3common.TransactionData) {
	param := steloParam(chainId, transactionData)
	apikey := "eDVo/n4nvQ0izEy/lDvw_AjD"

	postData, _ := json.Marshal(param)

	curl := fmt.Sprintf(`curl https://app.steloapi.com/api/v0/transaction?apiKey=%s -H 'Origin: https://docs.stelolabs.com' -H 'Referer: https://docs.stelolabs.com/' -H  'accept: application/json' -H 'content-type: application/json' --data-binary '%s'`, apikey, postData)

	fmt.Println("run the following command in terminal to send request to stelo")
	fmt.Println(curl)
}

func PrintAssetChangeBundleRequest(chainId uint, transactionsData []go2web3common.TransactionData) {
	param := make([]interface{}, len(transactionsData))
	for index, transactionData := range transactionsData {
		param[index] = steloParam(chainId, transactionData)
	}

	params := []interface{}{param}
	apikey := "eDVo/n4nvQ0izEy/lDvw_AjD"

	jsonrpc := map[string]interface{}{
		"id":      1,
		"jsonrpc": "2.0",
		"method":  "stelo_simulateAssetChangesBundle",
		"params":  params}

	postData, _ := json.Marshal(jsonrpc)

	curl := fmt.Sprintf(`curl https://app.steloapi.com/api/v0/transaction?apiKey=%s -H 'Origin: https://docs.stelolabs.com' -H 'Referer: https://docs.stelolabs.com/' -H  'accept: application/json' -H 'content-type: application/json' --data-binary '%s'`, apikey, postData)

	fmt.Println("run the following command in terminal to send request to stelo")
	fmt.Println(curl)
}
