package erc20Transaction

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"goland/go2web3/ethTransaction"
	"math/big"

	coreEntities "github.com/daoleno/uniswap-sdk-core/entities"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/metachris/eth-go-bindings/erc20"
	"golang.org/x/crypto/sha3"
)

func USDT(chainId uint) *coreEntities.Token {
	if chainId == 137 {
		return coreEntities.NewToken(chainId, common.HexToAddress("0xc2132D05D31c914a87C6611C10748AEb04B58e8F"), 6, "USDT", "USD Tether")
	} else if chainId == 1 {
		return coreEntities.NewToken(chainId, common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"), 6, "USDT", "USD Tether")
	} else {
		panic(errors.New(fmt.Sprintf("USDT undefined for %d", chainId)))
	}
}
func USDC(chainId uint) *coreEntities.Token {
	if chainId == 137 {
		return coreEntities.NewToken(chainId, common.HexToAddress("0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174"), 6, "USDC", "USD Coin")
	} else if chainId == 1 {
		return coreEntities.NewToken(chainId, common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"), 6, "USDC", "USD Coin")
	} else {
		panic(errors.New(fmt.Sprintf("USDT undefined for %d", chainId)))
	}
}
func UNI(chainId uint) *coreEntities.Token {
	if chainId == 137 {
		return coreEntities.NewToken(chainId, common.HexToAddress("0xb33EaAd8d922B1083446DC23f610c2567fB5180f"), 18, "UNI", "UNISwap Token")
	} else if chainId == 1 {
		return coreEntities.NewToken(chainId, common.HexToAddress("0x1f9840a85d5aF5bf1D1762F925BDADdC4201F984"), 18, "UNI", "UNISwap Token")
	} else {
		panic(errors.New(fmt.Sprintf("UNI undefined for %d", chainId)))
	}
}
func WrappedMatic(chainId uint) *coreEntities.Token {
	if chainId == 137 {
		return coreEntities.NewToken(chainId, common.HexToAddress("0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270"), 18, "WMATIC", "Wrapped Matic")
	} else if chainId == 1 {
		return coreEntities.NewToken(chainId, common.HexToAddress("0x7c9f4C87d911613Fe9ca58b579f737911AAD2D43"), 18, "WMATIC", "Wrapped Matic (Wormhole)")
	} else {
		panic(errors.New(fmt.Sprintf("WrappedMatic undefined for %d", chainId)))
	}
}
func WrappedEther(chainId uint) *coreEntities.Token {
	if chainId == 137 {
		return coreEntities.NewToken(chainId, common.HexToAddress("0x7ceB23fD6bC0adD59E62ac25578270cFf1b9f619"), 18, "WETH", "Wrapped Ether")
	} else if chainId == 1 {
		return coreEntities.NewToken(chainId, common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"), 18, "WETH", "Wrapped Ether")
	} else if chainId == 5 {
		return coreEntities.NewToken(chainId, common.HexToAddress("0xB4FBF271143F4FBf7B91A5ded31805e42b2208d6"), 18, "WETH", "Wrapped Ether")
	} else {
		panic(errors.New(fmt.Sprintf("WrappedEther undefined for %d", chainId)))
	}
}

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
	amount float64) (error, []byte, common.Address, uint64, *big.Int) {

	token, err := erc20.NewErc20(contractAddress, client)
	if err != nil {
		return err, nil, common.Address{}, 0, nil
	}

	decimals, err := token.Decimals(&bind.CallOpts{})
	if err != nil {
		return err, nil, common.Address{}, 0, nil
	}

	bigAmount := ethTransaction.FloatToBigInt(amount, int(decimals))

	callData := transferCallData(toAddress, bigAmount)

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

	gasLimit := uint64(4 * 21000)

	return nil, callData, contractAddress, gasLimit, big.NewInt(0)
}

func Approve(client *ethclient.Client,
	privateKey *ecdsa.PrivateKey,
	contractAddress common.Address,
	toAddress common.Address,
	amount float64) (error, []byte, common.Address, uint64, *big.Int) {

	token, err := erc20.NewErc20(contractAddress, client)
	if err != nil {
		return err, nil, common.Address{}, 0, nil
	}

	decimals, err := token.Decimals(&bind.CallOpts{})
	if err != nil {
		return err, nil, common.Address{}, 0, nil
	}

	bigAmount := ethTransaction.FloatToBigInt(amount, int(decimals))

	callData := approveTransferCallData(toAddress, bigAmount)

	gasLimit := uint64(4 * 21000)

	return nil, callData, contractAddress, gasLimit, big.NewInt(0)
}
