package ethTransaction

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GoerliClient() (*ethclient.Client, error) {
	return ethclient.Dial("https://goerli.infura.io/v3/590d5bfc5f2840f0aa08547c66c83753")
}
func MainnetClient() (*ethclient.Client, error) {
	return ethclient.Dial("https://mainnet.infura.io/v3/590d5bfc5f2840f0aa08547c66c83753")
}
func PolygonClient() (*ethclient.Client, error) {
	return ethclient.Dial("https://polygon-rpc.com/")
}
func PolygonMumbaiClient() (*ethclient.Client, error) {
	return ethclient.Dial("https://rpc-mumbai.maticvigil.com/")
}
func OwnEdgeTestnetClient() (*ethclient.Client, error) {
	return ethclient.Dial("http://0.0.0.0:10012")
}

func GetTheKey() (*ecdsa.PrivateKey, error) {
	// 0x73C17c616B918c73d12EE745Ebe513DdBC0faFaC address
	return crypto.HexToECDSA("0b7e449896aacdec13a0bf13879ba5db505bd19b80d18106cebf3bfc839871f3")
}

func GenerateKey() (*ecdsa.PrivateKey, error) {
	return crypto.GenerateKey()
}

func FloatToBigInt(amount float64, decimals int) *big.Int {
	fAmount := new(big.Float).SetFloat64(amount)
	fi, _ := new(big.Float).Mul(fAmount, big.NewFloat(math.Pow10(decimals))).Int(nil)
	return fi
}

func BigIntToFloat(amount *big.Int, decimals int) float64 {
	fAmount := new(big.Float).SetInt(amount)
	fi, _ := new(big.Float).Mul(fAmount, big.NewFloat(math.Pow10(-decimals))).Float64()
	return fi
}

func Execute(client *ethclient.Client,
	privateKey *ecdsa.PrivateKey,
	address common.Address,
	gasLimit uint64,
	callData []byte,
	gasPrice *big.Int,
	ethValue *big.Int) error {

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	//fromAddress := common.HexToAddress("0x7b596aBf5C33B9b57BD7E9679B5131E9D5378b05")
	currentNonce, err := client.PendingNonceAt(context.Background(), fromAddress)

	fmt.Println(fromAddress)
	fmt.Println(currentNonce)

	transaction := types.NewTransaction(currentNonce, address, ethValue, gasLimit, gasPrice, callData)

	// signereth := types.NewEIP155Signer(big.NewInt(1)) // ethereum mainnet
	signereth := types.NewEIP155Signer(big.NewInt(137)) // polygon mainnet
	// signereth := types.NewEIP155Signer(big.NewInt(80001))
	hash := signereth.Hash(transaction)
	fmt.Println(hash)

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		panic(err)
	}

	transaction, err = transaction.WithSignature(signereth, signature)
	if err != nil {
		panic(err)
	}

	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(transactionJSON))

	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
	if err != nil {
		panic(err)
	}

	publicKeyBytes := crypto.FromECDSAPub(&privateKey.PublicKey)
	matches := bytes.Equal(sigPublicKey, publicKeyBytes)
	fmt.Println(matches)

	//err = client.SendTransaction(context.Background(), transaction)

	return err
}
