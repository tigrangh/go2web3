package ethTransaction

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func EthereumGoerliClient() (*ethclient.Client, uint, error) {
	client, err := ethclient.Dial("https://goerli.infura.io/v3/590d5bfc5f2840f0aa08547c66c83753")
	return client, 5, err
}
func EthereumMainnetClient() (*ethclient.Client, uint, error) {
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/590d5bfc5f2840f0aa08547c66c83753")
	return client, 1, err
}
func PolygonMainnetClient() (*ethclient.Client, uint, error) {
	client, err := ethclient.Dial("https://polygon-rpc.com/")
	return client, 137, err
}
func PolygonMumbaiClient() (*ethclient.Client, uint, error) {
	client, err := ethclient.Dial("https://rpc-mumbai.maticvigil.com/")
	return client, 80001, err
}
func OwnEdgeTestnetClient() (*ethclient.Client, uint, error) {
	client, err := ethclient.Dial("http://0.0.0.0:10012")
	return client, 0, err
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

func signTransaction(client *ethclient.Client,
	privateKey *ecdsa.PrivateKey,
	address common.Address,
	gasLimit uint64,
	callData []byte,
	chainId uint,
	etherValue *big.Int) (error, string, common.Hash, *big.Int, *types.Transaction) {

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	currentNonce, err := client.PendingNonceAt(context.Background(), fromAddress)

	gasPrice, err := client.SuggestGasPrice(context.Background())

	if nil != err {
		return err, "", common.Hash{}, nil, nil
	}

	transaction := types.NewTransaction(currentNonce, address, etherValue, gasLimit, gasPrice, callData)

	signereth := types.NewEIP155Signer(big.NewInt(int64(chainId)))
	hash := signereth.Hash(transaction)

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return err, "", common.Hash{}, nil, nil
	}

	transaction, err = transaction.WithSignature(signereth, signature)
	if err != nil {
		return err, "", common.Hash{}, nil, nil
	}

	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		return err, "", common.Hash{}, nil, nil
	}

	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
	if err != nil {
		return err, "", common.Hash{}, nil, nil
	}

	publicKeyBytes := crypto.FromECDSAPub(&privateKey.PublicKey)
	if false == bytes.Equal(sigPublicKey, publicKeyBytes) {
		return errors.New("how can signature be incorrect?"), "", common.Hash{}, nil, nil
	}

	return err, string(transactionJSON), hash, gasPrice, transaction
}

func SignTransaction(client *ethclient.Client,
	privateKey *ecdsa.PrivateKey,
	address common.Address,
	gasLimit uint64,
	callData []byte,
	chainId uint,
	etherValue *big.Int) (error, string, common.Hash, *big.Int) {

	err, transactionJSON, hash, gasPrice, _ := signTransaction(client, privateKey, address, gasLimit, callData, chainId, etherValue)
	return err, transactionJSON, hash, gasPrice
}

func SignTransactionAndExecute(client *ethclient.Client,
	privateKey *ecdsa.PrivateKey,
	address common.Address,
	gasLimit uint64,
	callData []byte,
	chainId uint,
	etherValue *big.Int) (error, string, common.Hash, *big.Int) {

	err, transactionJSON, hash, gasPrice, transaction := signTransaction(client, privateKey, address, gasLimit, callData, chainId, etherValue)
	if err != nil {
		return err, "", common.Hash{}, nil
	}

	err = client.SendTransaction(context.Background(), transaction)
	if err != nil {
		return err, "", common.Hash{}, nil
	}

	return err, transactionJSON, hash, gasPrice
}
