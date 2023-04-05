package erc20Transaction

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"goland/go2web3/ethTransaction"

	coreEntities "github.com/daoleno/uniswap-sdk-core/entities"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/metachris/eth-go-bindings/erc20"
	"golang.org/x/crypto/sha3"
)

var (
	// USDC = coreEntities.NewToken(5, common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"), 6, "USDC", "USD Coin")
	// USDT = coreEntities.NewToken(5, common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"), 6, "USDT", "USD Tether")
	// DAI  = coreEntities.NewToken(1, common.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F"), 18, "DAI", "Dai Stablecoin")

	// Uni0 = coreEntities.NewToken(5, common.HexToAddress("0x280D4CEFC9178398B3F358fF0dd13b38f14a9f2A"), 18, "UNI0", "Uni0 token")
	// Uni1 = coreEntities.NewToken(5, common.HexToAddress("0x1225c6990Dd8c715a9CE327F2510BEF253909855"), 18, "UNI1", "Uni1 token")

	// UNI  = coreEntities.NewToken(5, common.HexToAddress("0x1f9840a85d5aF5bf1D1762F925BDADdC4201F984"), 18, "UNI", "Uniswap")
	// WETH = coreEntities.NewToken(5, common.HexToAddress("0xB4FBF271143F4FBf7B91A5ded31805e42b2208d6"), 18, "WETH", "Wrapped Ether")

	WMATIC = coreEntities.NewToken(137, common.HexToAddress("0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270"), 18, "Matic", "Matic Network(PolyGon)")
	AMP    = coreEntities.NewToken(137, common.HexToAddress("0x0621d647cecbFb64b79E44302c1933cB4f27054d"), 18, "AMP", "Amp")
	USDC   = coreEntities.NewToken(137, common.HexToAddress("0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174"), 6, "USDC", "USD Coin")
	USDT   = coreEntities.NewToken(137, common.HexToAddress("0xc2132D05D31c914a87C6611C10748AEb04B58e8F"), 6, "USDT", "USD Tether")
	AGIX   = coreEntities.NewToken(137, common.HexToAddress("0x190Eb8a183D22a4bdf278c6791b152228857c033"), 6, "AGIX", "SingularityNET Token")
	UNI    = coreEntities.NewToken(137, common.HexToAddress("0xb33EaAd8d922B1083446DC23f610c2567fB5180f"), 18, "UNI", "UNISwap Token")
)

func transferCallData(toAddress common.Address, amount *big.Int) []byte {
	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)

	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	return data
}

func approveTransferCallData(toAddress common.Address, amount *big.Int) []byte {
	transferFnSignature := []byte("approve(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)

	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	return data
}

func Transfer(client *ethclient.Client,
	privateKey *ecdsa.PrivateKey,
	contractAddress common.Address,
	toAddress common.Address,
	amount float64) error {

	token, err := erc20.NewErc20(contractAddress, client)
	if err != nil {
		return err
	}

	decimals, err := token.Decimals(&bind.CallOpts{})
	if err != nil {
		return err
	}

	bigAmount := ethTransaction.FloatToBigInt(amount, int(decimals))

	callData := transferCallData(toAddress, bigAmount)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if nil != err {
		return err
	}

	// fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	// gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
	// 	From:     fromAddress,
	// 	To:       &contractAddress,
	// 	GasPrice: gasPrice,
	// 	Value:    bigAmount,
	// 	Data:     callData,
	// })
	// if nil != err {
	// 	return err
	// }

	gasLimit := uint64(4 * 21000) // works with this on Polygon

	return ethTransaction.Execute(client, privateKey, contractAddress, gasLimit, callData, gasPrice, big.NewInt(0))
}

func ApproveTransfer(client *ethclient.Client,
	privateKey *ecdsa.PrivateKey,
	contractAddress common.Address,
	toAddress common.Address,
	amount float64) error {

	token, err := erc20.NewErc20(contractAddress, client)
	if err != nil {
		return err
	}

	decimals, err := token.Decimals(&bind.CallOpts{})
	if err != nil {
		return err
	}

	bigAmount := ethTransaction.FloatToBigInt(amount, int(decimals))

	callData := approveTransferCallData(toAddress, bigAmount)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if nil != err {
		return err
	}

	// fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	// gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
	// 	From:     fromAddress,
	// 	To:       &contractAddress,
	// 	GasPrice: gasPrice,
	// 	Value:    bigAmount,
	// 	Data:     callData,
	// })
	// if nil != err {
	// 	return err
	// }

	gasLimit := uint64(4 * 21000) // works with this on Polygon

	return ethTransaction.Execute(client, privateKey, contractAddress, gasLimit, callData, gasPrice, big.NewInt(0))
}
