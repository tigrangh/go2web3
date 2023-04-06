package main

import (
	"fmt"
	"goland/go2web3/alchemy"
	"goland/go2web3/erc20Transaction"
	"goland/go2web3/ethTransaction"
	"goland/go2web3/uni"

	"github.com/ethereum/go-ethereum/crypto"
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

	// err = ethTransaction.TransferEther(client, privateKey, toAddress, 0.05)
	// err = erc20Transaction.TransferERC20(client, privateKey, erc20Transaction.USDC.Address, toAddress, 0.05)

	// err = erc20Transaction.ApproveTransfer(client,
	// 	privateKey,
	// 	erc20Transaction.USDC.Address,
	// 	common.HexToAddress("0xE592427A0AEce92De3Edee1F18E0157C05861564"),
	// 	10)

	wmatic := erc20Transaction.WrappedMatic(chainId)
	usdc := erc20Transaction.USDC(chainId)

	//fmt.Println(chainId, usdc.Address, wmatic.Address)

	err, callData, interactAddress, poolAddress, gasLimit, etherValue := uni.Swap(client,
		crypto.PubkeyToAddress(privateKey.PublicKey),
		100, // 0.01%
		usdc,
		1,
		wmatic)

	if nil != err {
		panic(err)
	}

	fmt.Printf("Pool Address: %s\n\n", poolAddress)

	err, transactionJSON, hash, gasPrice := ethTransaction.SignTransaction(client, privateKey, interactAddress, gasLimit, callData, chainId, etherValue)

	if nil != err {
		panic(err)
	}

	fmt.Printf("JSON transaction: %s\n\n", transactionJSON)
	fmt.Printf("Transaction hash: %s\n\n", hash)

	alchemy.PrintAssetChangeRequest(chainId, interactAddress, crypto.PubkeyToAddress(privateKey.PublicKey), etherValue, callData, gasLimit, gasPrice)
}
