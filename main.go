package main

import (
	"fmt"
	"goland/go2web3/alchemy"
	"goland/go2web3/erc20Transaction"
	"goland/go2web3/ethTransaction"

	"github.com/ethereum/go-ethereum/common"
)

func main() {

	client, chainId, err := ethTransaction.PolygonMainnetClient()
	if nil != err {
		panic(err)
	}

	// uni.GetRate(3000 /*0.3%*/, uni.UNI, 1, uni.USDT)

	// uni.GetRate(3000 /*0.3%*/, uni.USDC, 1, uni.USDT)
	// uni.GetRate(500 /*0.05%*/, uni.USDC, 1, uni.USDT)
	// uni.GetRate(100 /*0.01%*/, uni.USDC, 1, uni.USDT)

	// uni.GetRate(3000 /*0.3%*/, uni.USDT, 1, uni.USDC)
	// uni.GetRate(500 /*0.05%*/, uni.USDT, 1, uni.USDC)
	// uni.GetRate(100 /*0.01%*/, uni.USDT, 1, uni.USDC)

	//privateKey, err := ethTransaction.GenerateKey()
	privateKey, err := ethTransaction.GetTheKey()
	if nil != err {
		panic(err)
	}

	// uni.GetRate(client, 100 /*0.01%*/, erc20Transaction.USDC, 1, erc20Transaction.USDT)

	// uni.GetRate(client, 100 /*0.01%*/, erc20Transaction.UNI, 0.00005821, erc20Transaction.WETH)
	// uni.GetRate(client, 100 /*0.01%*/, erc20Transaction.WETH, 0.0000957675, erc20Transaction.UNI)
	// err = uni.GetRate(client, 100, erc20Transaction.USDC, 0.1, erc20Transaction.WMATIC)
	// toAddress := common.HexToAddress("0x1298bF10baa546A332D9c675d390f79Cf375227C")

	//wmatic := erc20Transaction.WrappedMatic(chainId)
	usdc := erc20Transaction.USDC(chainId)
	// gtc := erc20Transaction.GeghamToken(chainId)

	// fmt.Println(chainId, usdc.Address, wmatic.Address)
	// fmt.Println(crypto.PubkeyToAddress(privateKey.PublicKey))
	// fmt.Println(fmt.Println(hex.EncodeToString(crypto.FromECDSA(privateKey))))

	// err, callData, interactAddress, gasLimit, etherValue := erc20Transaction.Approve(client,
	// 	wmatic.Address,
	// 	common.HexToAddress("0xE592427A0AEce92De3Edee1F18E0157C05861564"),
	// 	4)

	err, callData, interactAddress, gasLimit, etherValue := erc20Transaction.Transfer(client,
		usdc.Address,
		common.HexToAddress("0x73c17c616b918c73d12ee745ebe513ddbc0fafac"),
		90)

	// err, callData, interactAddress, poolAddress, gasLimit, etherValue := uni.Swap(client,
	// 	crypto.PubkeyToAddress(privateKey.PublicKey),
	// 	100, // 0.01%
	// 	usdc,
	// 	1,
	// 	wmatic)

	if nil != err {
		panic(err)
	}

	//fmt.Printf("Pool Address: %s\n\n", poolAddress)

	err, transactionJSON, hash, _, transactionData := ethTransaction.SignTransaction(client,
		privateKey,
		interactAddress,
		gasLimit,
		callData,
		chainId,
		etherValue)

	// err, transactionJSON, hash, _, transactionData := ethTransaction.SignJSONTransaction(client,
	// 	privateKey,
	// 	`{"type":"0x0","nonce":"0x22","gasPrice":"0x25da7ae700","maxPriorityFeePerGas":null,"maxFeePerGas":null,"gas":"0x4ce78","value":"0x0","input":"0x414bf3890000000000000000000000000d500b1d8e8ef31e21c99d1db9a6444d3adf12700000000000000000000000002791bca1f2de4661ed88a30c99a7a9449aa841740000000000000000000000000000000000000000000000000000000000000064000000000000000000000000a7b1d2f8dcb87216f4876b8cd3828ae6b48e4d6e000000000000000000000000000000000000000000000000000000006430088600000000000000000000000000000000000000000000000030f32e32becc7a00000000000000000000000000000000000000000000000000000000000035e4f20000000000000000000000000000000000000000000000000000000000000000","v":"0x0","r":"0x0","s":"0x0","to":"0xe592427a0aece92de3edee1f18e0157c05861564","hash":"0x2ab217de2a0215614f9eb223a31a381346abafb2c999a0226778c22958c0865a"}`,
	// 	chainId)

	if nil != err {
		panic(err)
	}

	fmt.Printf("JSON transaction: %s\n\n", transactionJSON)
	fmt.Printf("Transaction hash: %s\n\n", hash)

	alchemy.PrintAssetChangeRequest(chainId, transactionData)

	// // alchemy.PrintAssetChangeRequest(chainId, interactAddress, crypto.PubkeyToAddress(privateKey.PublicKey), etherValue, callData, gasLimit, gasPrice)
	// alchemy.PrintAssetChangeRequest(chainId,
	// 	*transaction.To(),
	// 	crypto.PubkeyToAddress(privateKey.PublicKey),
	// 	transaction.Value(),
	// 	transaction.Data(),
	// 	transaction.Gas(),
	// 	gasPrice)

}
