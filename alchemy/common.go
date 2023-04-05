package alchemy

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func PrintAssetChangeRequest(interactAddress common.Address,
	fromAddress common.Address,
	etherValue *big.Int,
	callData []byte,
	gasLimit uint64,
	gasPrice *big.Int) {

	post := fmt.Sprintf(`{"id": 1, "jsonrpc": "2.0", "method":"alchemy_simulateAssetChanges","params": [{"from": "%s", "to": "%s", "value": "0x%x", "data": "0x%x", "gas": "0x%x", "gasPrice": "0x%x"}]}`, fromAddress, interactAddress, etherValue, callData, gasLimit, gasPrice)

	curl := fmt.Sprintf(`curl https://polygon-mainnet.g.alchemy.com/v2/docs-demo -H 'Origin: https://docs.alchemy.com' -H 'accept: application/json' -H 'content-type: application/json' --data-binary '%s'`, post)

	fmt.Printf("copy paste to send request to alchemy\n\n")
	fmt.Println(curl)
}
