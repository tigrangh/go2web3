package ethTransaction

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"goland/go2web3/go2web3common"
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

func prepareTransactionDataByJSON(transactionJSON string,
	fromAddress common.Address) go2web3common.TransactionData {

	transaction := new(types.Transaction)
	json.Unmarshal([]byte(transactionJSON), transaction)

	return go2web3common.TransactionData{
		FromAddress:     fromAddress,
		InteractAddress: *transaction.To(),
		EtherValue:      transaction.Value(),
		CallData:        transaction.Data(),
		GasLimit:        transaction.Gas(),
		GasPrice:        transaction.GasPrice()}
}

func prepareTransactionByData(client *ethclient.Client,
	transactionData go2web3common.TransactionData) (error, *types.Transaction, go2web3common.TransactionData) {

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if nil != err {
		return err, nil, go2web3common.TransactionData{}
	}

	currentNonce, err := client.PendingNonceAt(context.Background(), transactionData.FromAddress)
	if nil != err {
		return err, nil, go2web3common.TransactionData{}
	}

	transactionData.GasPrice = gasPrice

	return nil, types.NewTransaction(currentNonce,
		transactionData.InteractAddress,
		transactionData.EtherValue,
		transactionData.GasLimit,
		transactionData.GasPrice,
		transactionData.CallData), transactionData
}

func signTransaction(client *ethclient.Client,
	privateKey *ecdsa.PrivateKey,
	transaction *types.Transaction,
	chainId uint) (error, string, common.Hash, *types.Transaction) {

	signereth := types.NewEIP155Signer(big.NewInt(int64(chainId)))
	hash := signereth.Hash(transaction)

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return err, "", common.Hash{}, nil
	}

	transaction, err = transaction.WithSignature(signereth, signature)
	if err != nil {
		return err, "", common.Hash{}, nil
	}

	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		return err, "", common.Hash{}, nil
	}

	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
	if err != nil {
		return err, "", common.Hash{}, nil
	}

	publicKeyBytes := crypto.FromECDSAPub(&privateKey.PublicKey)
	if false == bytes.Equal(sigPublicKey, publicKeyBytes) {
		return errors.New("how can signature be incorrect?"), "", common.Hash{}, nil
	}

	return err, string(transactionJSON), hash, transaction
}

func SignTransaction(client *ethclient.Client,
	privateKey *ecdsa.PrivateKey,
	interactAddress common.Address,
	gasLimit uint64,
	callData []byte,
	chainId uint,
	etherValue *big.Int) (error, string, common.Hash, *types.Transaction, go2web3common.TransactionData) {

	transactionData := go2web3common.TransactionData{
		FromAddress:     crypto.PubkeyToAddress(privateKey.PublicKey),
		InteractAddress: interactAddress,
		EtherValue:      etherValue,
		CallData:        callData,
		GasLimit:        gasLimit,
		GasPrice:        big.NewInt(0)}

	err, transaction, transactionData := prepareTransactionByData(client, transactionData)
	if err != nil {
		return err, "", common.Hash{}, nil, go2web3common.TransactionData{}
	}

	err, transactionJSON, hash, transaction := signTransaction(client, privateKey, transaction, chainId)
	return err, transactionJSON, hash, transaction, transactionData
}

func SignTransactionAndExecute(client *ethclient.Client,
	privateKey *ecdsa.PrivateKey,
	interactAddress common.Address,
	gasLimit uint64,
	callData []byte,
	chainId uint,
	etherValue *big.Int) (error, string, common.Hash, *types.Transaction, go2web3common.TransactionData) {

	transactionData := go2web3common.TransactionData{
		FromAddress:     crypto.PubkeyToAddress(privateKey.PublicKey),
		InteractAddress: interactAddress,
		EtherValue:      etherValue,
		CallData:        callData,
		GasLimit:        gasLimit,
		GasPrice:        big.NewInt(0)}

	err, transaction, transactionData := prepareTransactionByData(client, transactionData)
	if err != nil {
		return err, "", common.Hash{}, nil, go2web3common.TransactionData{}
	}

	err, transactionJSON, hash, transaction := signTransaction(client, privateKey, transaction, chainId)
	if err != nil {
		return err, "", common.Hash{}, nil, go2web3common.TransactionData{}
	}

	err = client.SendTransaction(context.Background(), transaction)
	if err != nil {
		return err, "", common.Hash{}, nil, go2web3common.TransactionData{}
	}

	return err, transactionJSON, hash, transaction, transactionData
}

func SignJSONTransaction(client *ethclient.Client,
	privateKey *ecdsa.PrivateKey,
	transactionJSON string,
	chainId uint) (error, string, common.Hash, *types.Transaction, go2web3common.TransactionData) {

	transactionData := prepareTransactionDataByJSON(transactionJSON, crypto.PubkeyToAddress(privateKey.PublicKey))

	err, transaction, transactionData := prepareTransactionByData(client, transactionData)
	if err != nil {
		return err, "", common.Hash{}, nil, go2web3common.TransactionData{}
	}

	err, transactionJSON, hash, transaction := signTransaction(client, privateKey, transaction, chainId)
	return err, transactionJSON, hash, transaction, transactionData
}

func SignJSONTransactionAndExecute(client *ethclient.Client,
	privateKey *ecdsa.PrivateKey,
	transactionJSON string,
	chainId uint) (error, string, common.Hash, *types.Transaction, go2web3common.TransactionData) {

	transactionData := prepareTransactionDataByJSON(transactionJSON, crypto.PubkeyToAddress(privateKey.PublicKey))

	err, transaction, transactionData := prepareTransactionByData(client, transactionData)
	if err != nil {
		return err, "", common.Hash{}, nil, go2web3common.TransactionData{}
	}

	err, transactionJSON, hash, transaction := signTransaction(client, privateKey, transaction, chainId)
	if err != nil {
		return err, "", common.Hash{}, nil, go2web3common.TransactionData{}
	}

	err = client.SendTransaction(context.Background(), transaction)
	if err != nil {
		return err, "", common.Hash{}, nil, go2web3common.TransactionData{}
	}

	return err, transactionJSON, hash, transaction, transactionData
}
