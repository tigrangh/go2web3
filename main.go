package main

import (
	"goland/go2web3/erc20Transaction"
	"goland/go2web3/ethTransaction"
	"goland/go2web3/uni"
)

//"github.com/ethereum/go-ethereum/common"

func main() {

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

	client, err := ethTransaction.PolygonClient()
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

	err = uni.Swap(client,
		privateKey,
		100, // 0.01%
		erc20Transaction.USDC,
		1,
		erc20Transaction.WMATIC)

	if nil != err {
		panic(err)
	}
}
