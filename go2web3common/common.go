package go2web3common

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type TransactionData struct {
	FromAddress     common.Address
	InteractAddress common.Address
	EtherValue      *big.Int
	CallData        []byte
	GasLimit        uint64
	GasPrice        *big.Int
}
