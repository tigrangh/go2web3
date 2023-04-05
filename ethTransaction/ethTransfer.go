package ethTransaction

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TransferEther(client *ethclient.Client,
	privateKey *ecdsa.PrivateKey,
	toAddress common.Address,
	amount float64) ([]byte, common.Address, uint64, *big.Int) {

	gasLimit := uint64(21000)
	var callData []byte

	bigAmount := FloatToBigInt(amount, 18)

	return callData, toAddress, gasLimit, bigAmount
}
