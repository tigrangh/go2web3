package alchemy

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type RequestBody struct {
	//ChainId         uint
	InteractAddress common.Address
	FromAddress     common.Address
	EtherValue      *big.Int
	CallData        []byte
	GasLimit        uint64
	GasPrice        *big.Int
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

func PrintAssetChangeRequest(chainId uint, obj RequestBody) {
	endpointPart := DecideAlchemyEndpoint(chainId)

	post := fmt.Sprintf(`{"id": 1, "jsonrpc": "2.0", "method":"alchemy_simulateAssetChanges","params": [{"from": "%s", "to": "%s", "value": "0x%x", "data": "0x%x", "gas": "0x%x", "gasPrice": "0x%x"}]}`,
		obj.FromAddress, obj.InteractAddress, obj.EtherValue, obj.CallData, obj.GasLimit, obj.GasPrice)

	curl := fmt.Sprintf(`curl https://%s.g.alchemy.com/v2/docs-demo -H 'Origin: https://docs.alchemy.com' -H 'accept: application/json' -H 'content-type: application/json' --data-binary '%s'`, endpointPart, post)

	fmt.Println("run the following command in terminal to send request to alchemy")
	fmt.Println(curl)
}

func PrintAssetChangeBundleRequest(chainId uint, obj []RequestBody) {
	//var result string
	var post []string
	post = append(post, fmt.Sprintf(`{"id": 1, "jsonrpc": "2.0", "method":"alchemy_simulateAssetChanges"}`))

	endpointPart := DecideAlchemyEndpoint(chainId)
	for i := 0; i < len(obj); i++ {
		// post = fmt.Sprintf(`{"id": 1, "jsonrpc": "2.0", "method":"alchemy_simulateAssetChanges","params": [[{"from": "%s", "to": "%s", "value": "0x%x", "data": "0x%x", "gas": "0x%x", "gasPrice": "0x%x"}]]}`,
		// 	obj[i].FromAddress, obj[i].InteractAddress, obj[i].EtherValue, obj[i].CallData, obj[i].GasLimit, obj[i].GasPrice)
		// result += post

		post = append(post, fmt.Sprintf(`{"params": [{"from": "%s", "to": "%s", "value": "0x%x", "data": "0x%x", "gas": "0x%x", "gasPrice": "0x%x"}]}`,
			obj[i].FromAddress, obj[i].InteractAddress, obj[i].EtherValue, obj[i].CallData, obj[i].GasLimit, obj[i].GasPrice))
	}

	curl := fmt.Sprintf(`curl https://%s.g.alchemy.com/v2/docs-demo -H 'Origin: https://docs.alchemy.com' -H 'accept: application/json' -H 'content-type: application/json' --data-binary '%s'`, endpointPart, post)

	fmt.Println("run the following command in terminal to send request to alchemy")
	fmt.Println(curl)

}
