package uni

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"goland/go2web3/ethTransaction"
	"time"

	coreEntities "github.com/daoleno/uniswap-sdk-core/entities"
	"github.com/daoleno/uniswapv3-sdk/constants"
	"github.com/daoleno/uniswapv3-sdk/entities"
	"github.com/daoleno/uniswapv3-sdk/examples/contract"
	"github.com/daoleno/uniswapv3-sdk/examples/helper"
	"github.com/daoleno/uniswapv3-sdk/periphery"
	sdkutils "github.com/daoleno/uniswapv3-sdk/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Swap(client *ethclient.Client,
	privateKey *ecdsa.PrivateKey,
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

	gasPrice, err := client.SuggestGasPrice(context.Background())
	//gasPrice := ethTransaction.FloatToBigInt(1.879053, 12) // this is a sample value from a successful tx example

	if nil != err {
		return err
	}

	fmt.Println("pool address: ", poolAddress)
	uniswapV3Pool, err := contract.NewUniswapv3Pool(poolAddress, client)
	if err != nil {
		return err
	}

	liquidity, err := uniswapV3Pool.Liquidity(nil)
	if err != nil {
		return err
	}

	slot0, err := uniswapV3Pool.Slot0(nil)
	if err != nil {
		return err
	}

	pooltick, err := uniswapV3Pool.Ticks(nil, big.NewInt(0))
	if err != nil {
		return err
	}

	feeAmount := constants.FeeAmount(uint64(swapFee))
	ticks := []entities.Tick{
		{
			Index: entities.NearestUsableTick(sdkutils.MinTick,
				constants.TickSpacings[feeAmount]),
			LiquidityNet:   pooltick.LiquidityNet,
			LiquidityGross: pooltick.LiquidityGross,
		},
		{
			Index: entities.NearestUsableTick(sdkutils.MaxTick,
				constants.TickSpacings[feeAmount]),
			LiquidityNet:   pooltick.LiquidityNet,
			LiquidityGross: pooltick.LiquidityGross,
		},
	}

	tickListDataProvider, err := entities.NewTickListDataProvider(ticks, constants.TickSpacings[feeAmount])
	if err != nil {
		return err
	}

	poolEntity, err := entities.NewPool(tokenIn,
		tokenOut,
		constants.FeeAmount(uint64(swapFee)),
		slot0.SqrtPriceX96,
		liquidity,
		int(slot0.Tick.Int64()),
		tickListDataProvider)

	if err != nil {
		return err
	}

	routeEntity, err := entities.NewRoute([]*entities.Pool{poolEntity}, tokenIn, tokenOut)
	if err != nil {
		return err
	}

	bigAmountIn := ethTransaction.FloatToBigInt(amountIn, int(tokenIn.Decimals()))

	//10%
	slippageTolerance := coreEntities.NewPercent(big.NewInt(1), big.NewInt(10))

	currentTimePlus5Minutes := time.Now().Add(time.Minute * time.Duration(5)).Unix()
	deadlineCurrentTimePlus5Minutes := big.NewInt(currentTimePlus5Minutes)

	tradeEntity, err := entities.FromRoute(routeEntity, coreEntities.FromRawAmount(tokenIn, bigAmountIn), coreEntities.ExactInput)
	if err != nil {
		return err
	}

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	//fromAddress := common.HexToAddress("0x7b596aBf5C33B9b57BD7E9679B5131E9D5378b05")
	//fromAddress := common.HexToAddress("0x8B26320912935111300DdAeeC15EA9a182FF6F1a")

	params, err := periphery.SwapCallParameters([]*entities.Trade{tradeEntity}, &periphery.SwapOptions{
		SlippageTolerance: slippageTolerance,
		Recipient:         fromAddress,
		Deadline:          deadlineCurrentTimePlus5Minutes,
	})
	if err != nil {
		return err
	}

	swapRouterAddress := common.HexToAddress(helper.ContractV3SwapRouterV1)

	// gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
	// 	From:     fromAddress,
	// 	To:       &swapRouterAddress,
	// 	GasPrice: gasPrice,
	// 	Value:    bigAmountIn,
	// 	Data:     params.Calldata,
	// })
	// if err != nil {
	// 	return err
	// }

	gasLimit := uint64(15 * 21000)

	callData := params.Calldata

	return ethTransaction.Execute(client, privateKey, swapRouterAddress, gasLimit, callData, gasPrice, big.NewInt(0))
}
