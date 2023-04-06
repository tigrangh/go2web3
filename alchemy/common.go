package alchemy

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func PrintAssetChangeRequest(chainId uint,
	interactAddress common.Address,
	fromAddress common.Address,
	etherValue *big.Int,
	callData []byte,
	gasLimit uint64,
	gasPrice *big.Int) {

	endpointPart := ""
	if chainId == 1 {
		endpointPart = "eth-mainnet"
	} else if chainId == 137 {
		endpointPart = "polygon-mainnet"
	} else if chainId == 5 {
		endpointPart = "eth-goerly"
	} else if chainId == 80001 {
		endpointPart = "polygon-mumbai"
	} else {
		panic(errors.New(fmt.Sprintf("alchemy endpoint undefined for %d", chainId)))
	}

	post := fmt.Sprintf(`{"id": 1, "jsonrpc": "2.0", "method":"alchemy_simulateAssetChanges","params": [{"from": "%s", "to": "%s", "value": "0x%x", "data": "0x%x", "gas": "0x%x", "gasPrice": "0x%x"}]}`, fromAddress, interactAddress, etherValue, callData, gasLimit, gasPrice)

	curl := fmt.Sprintf(`curl https://%s.g.alchemy.com/v2/docs-demo -H 'Origin: https://docs.alchemy.com' -H 'accept: application/json' -H 'content-type: application/json' --data-binary '%s'`, endpointPart, post)

	fmt.Println("run the following command in terminal to send request to alchemy")
	fmt.Println(curl)
}
