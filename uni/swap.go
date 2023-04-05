package uni

import (
	"errors"
	"goland/go2web3/ethTransaction"
	"math/big"
	"time"

	coreEntities "github.com/daoleno/uniswap-sdk-core/entities"
	"github.com/daoleno/uniswapv3-sdk/constants"
	"github.com/daoleno/uniswapv3-sdk/entities"
	"github.com/daoleno/uniswapv3-sdk/examples/contract"
	"github.com/daoleno/uniswapv3-sdk/examples/helper"
	"github.com/daoleno/uniswapv3-sdk/periphery"
	sdkutils "github.com/daoleno/uniswapv3-sdk/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func Swap(client *ethclient.Client,
	recipient common.Address,
	swapFee int64,
	tokenIn *coreEntities.Token,
	amountIn float64,
	tokenOut *coreEntities.Token) (error, []byte, common.Address, common.Address, uint64, *big.Int) {

	uniswapv3Factory, err := contract.NewUniswapv3Factory(common.HexToAddress(helper.ContractV3Factory), client)
	if err != nil {
		return err, nil, common.Address{}, common.Address{}, 0, nil
	}
	uniswapv3FactoryRaw := &contract.Uniswapv3FactoryRaw{Contract: uniswapv3Factory}

	swapFeeBigInt := big.NewInt(swapFee)

	var outGetPool []interface{}
	err = uniswapv3FactoryRaw.Call(nil, &outGetPool, "getPool", tokenIn.Address, tokenOut.Address, swapFeeBigInt)
	if err != nil {
		return err, nil, common.Address{}, common.Address{}, 0, nil
	}
	if 0 == len(outGetPool) {
		return errors.New("no pool is found"), nil, common.Address{}, common.Address{}, 0, nil
	}
	poolAddress := outGetPool[0].(common.Address)
	if poolAddress == (common.Address{}) {
		return errors.New("no pool is found"), nil, common.Address{}, common.Address{}, 0, nil
	}

	uniswapV3Pool, err := contract.NewUniswapv3Pool(poolAddress, client)
	if err != nil {
		return err, nil, common.Address{}, common.Address{}, 0, nil
	}

	liquidity, err := uniswapV3Pool.Liquidity(nil)
	if err != nil {
		return err, nil, common.Address{}, common.Address{}, 0, nil
	}

	slot0, err := uniswapV3Pool.Slot0(nil)
	if err != nil {
		return err, nil, common.Address{}, common.Address{}, 0, nil
	}

	pooltick, err := uniswapV3Pool.Ticks(nil, big.NewInt(0))
	if err != nil {
		return err, nil, common.Address{}, common.Address{}, 0, nil
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
		return err, nil, common.Address{}, common.Address{}, 0, nil
	}

	poolEntity, err := entities.NewPool(tokenIn,
		tokenOut,
		constants.FeeAmount(uint64(swapFee)),
		slot0.SqrtPriceX96,
		liquidity,
		int(slot0.Tick.Int64()),
		tickListDataProvider)

	if err != nil {
		return err, nil, common.Address{}, common.Address{}, 0, nil
	}

	routeEntity, err := entities.NewRoute([]*entities.Pool{poolEntity}, tokenIn, tokenOut)
	if err != nil {
		return err, nil, common.Address{}, common.Address{}, 0, nil
	}

	bigAmountIn := ethTransaction.FloatToBigInt(amountIn, int(tokenIn.Decimals()))

	//0.1%
	slippageTolerance := coreEntities.NewPercent(big.NewInt(1), big.NewInt(10))

	currentTimePlus5Minutes := time.Now().Add(time.Minute * time.Duration(5)).Unix()
	deadlineCurrentTimePlus5Minutes := big.NewInt(currentTimePlus5Minutes)

	tradeEntity, err := entities.FromRoute(routeEntity, coreEntities.FromRawAmount(tokenIn, bigAmountIn), coreEntities.ExactInput)
	if err != nil {
		return err, nil, common.Address{}, common.Address{}, 0, nil
	}

	params, err := periphery.SwapCallParameters([]*entities.Trade{tradeEntity}, &periphery.SwapOptions{
		SlippageTolerance: slippageTolerance,
		Recipient:         recipient,
		Deadline:          deadlineCurrentTimePlus5Minutes,
	})
	if err != nil {
		return err, nil, common.Address{}, common.Address{}, 0, nil
	}

	swapRouterAddress := common.HexToAddress(helper.ContractV3SwapRouterV1)

	gasLimit := uint64(15 * 21000)

	callData := params.Calldata

	return nil, callData, swapRouterAddress, poolAddress, gasLimit, big.NewInt(0)
}
