package ethTransaction

import (
	"context"
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TransferEther(client *ethclient.Client,
	privateKey *ecdsa.PrivateKey,
	toAddress common.Address,
	amount float64) error {

	gasLimit := uint64(21000)
	var callData []byte

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if nil != err {
		return err
	}

	bigAmount := FloatToBigInt(amount, 18)

	return Execute(client, privateKey, toAddress, gasLimit, callData, gasPrice, bigAmount)
}
