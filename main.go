package main

import (
	"fmt"
	"goland/go2web3/alchemy"
	"goland/go2web3/erc20Transaction"
	"goland/go2web3/ethTransaction"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {

	client, chainId, err := ethTransaction.PolygonMumbaiClient()
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

	// wmatic := erc20Transaction.WrappedMatic(chainId)
	// usdc := erc20Transaction.USDC(chainId)
	toAddress := common.HexToAddress("0x827C81e18a1b729Da4C0bf20c1Dedd0EC6dA2182")

	err, callData, interactAddress, gasLimit, etherValue := erc20Transaction.Transfer(client,
		privateKey, common.HexToAddress("0x18d223fe5d767b2db7586ad23d8134dd8990fa47"), toAddress, 100)

	//fmt.Println(chainId, usdc.Address, wmatic.Address)

	//fmt.Println(crypto.PubkeyToAddress(privateKey.PublicKey))

	//err, callData, interactAddress, poolAddress, gasLimit, etherValue := uni.Swap(client,
	// crypto.PubkeyToAddress(privateKey.PublicKey),
	// 100, // 0.01%
	// usdc,
	// 1,
	// wmatic)

	if nil != err {
		panic(err)
	}

	//fmt.Printf("Pool Address: %s\n\n", poolAddress)

	err, transactionJSON, hash, gasPrice := ethTransaction.SignTransaction(client, privateKey, interactAddress, gasLimit, callData, chainId, etherValue)

	if nil != err {
		panic(err)
	}

	fmt.Printf("JSON transaction: %s\n\n", transactionJSON)
	fmt.Printf("Transaction hash: %s\n\n", hash)

	ob := []alchemy.RequestBody{
		{
			InteractAddress: interactAddress,
			FromAddress:     crypto.PubkeyToAddress(privateKey.PublicKey),
			EtherValue:      etherValue,
			CallData:        callData,
			GasLimit:        gasLimit,
			GasPrice:        gasPrice,
		},
		{
			InteractAddress: interactAddress,
			FromAddress:     crypto.PubkeyToAddress(privateKey.PublicKey),
			EtherValue:      etherValue,
			CallData:        callData,
			GasLimit:        gasLimit,
			GasPrice:        gasPrice,
		},
		{
			InteractAddress: interactAddress,
			FromAddress:     crypto.PubkeyToAddress(privateKey.PublicKey),
			EtherValue:      etherValue,
			CallData:        callData,
			GasLimit:        gasLimit,
			GasPrice:        gasPrice,
		},
	}

	alchemy.PrintAssetChangeBundleRequest(chainId, ob) //chainId, interactAddress, crypto.PubkeyToAddress(privateKey.PublicKey), etherValue, callData, gasLimit, gasPrice
}
