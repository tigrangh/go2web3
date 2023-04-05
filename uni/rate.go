package uni

import (
	"errors"
	"fmt"
	"math/big"
	"goland/go2web3/ethTransaction"

	coreEntities "github.com/daoleno/uniswap-sdk-core/entities"
	"github.com/daoleno/uniswapv3-sdk/examples/contract"
	"github.com/daoleno/uniswapv3-sdk/examples/helper"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func GetRate(client *ethclient.Client,
	swapFee int64,
	tokenIn *coreEntities.Token,
	amountIn float64,
	tokenOut *coreEntities.Token) error {

	uniswapv3Factory, err := contract.NewUniswapv3Factory(common.HexToAddress(helper.ContractV3Factory), client)
	if err != nil {
		return err
	}
	uniswapv3FactoryRaw := &contract.Uniswapv3FactoryRaw{Contract: uniswapv3Factory}

	swapFeeBigInt := big.NewInt(swapFee)
	sqrtPriceLimitX96 := big.NewInt(0)

	var outGetPool []interface{}
	err = uniswapv3FactoryRaw.Call(nil, &outGetPool, "getPool", tokenIn.Address, tokenOut.Address, swapFeeBigInt)
	if err != nil {
		return err
	}
	if 0 == len(outGetPool) {
		return errors.New("no pool is found")
	}
	poolAddress := outGetPool[0].(common.Address)
	if poolAddress == (common.Address{}) {
		return errors.New("no pool is found")
	}

	fmt.Println("pool address: ", poolAddress)

	uniswapv3Quoter, err := contract.NewUniswapv3Quoter(common.HexToAddress(helper.ContractV3Quoter), client)
	if err != nil {
		return err
	}

	bigAmountIn := ethTransaction.FloatToBigInt(amountIn, int(tokenIn.Decimals()))

	var outQuote []interface{}
	uniswapv3QuoterRaw := &contract.Uniswapv3QuoterRaw{Contract: uniswapv3Quoter}
	err = uniswapv3QuoterRaw.Call(nil, &outQuote, "quoteExactInputSingle", tokenIn.Address, tokenOut.Address, swapFeeBigInt, bigAmountIn, sqrtPriceLimitX96)
	if err != nil {
		return err
	}

	fmt.Println("amountOut: ", ethTransaction.BigIntToFloat(outQuote[0].(*big.Int), int(tokenOut.Decimals())))

	return nil
}
