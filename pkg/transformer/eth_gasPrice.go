package transformer

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/revo"
)

// ProxyETHEstimateGas implements ETHProxy
type ProxyETHGasPrice struct {
	*revo.Revo
}

func (p *ProxyETHGasPrice) Method() string {
	return "eth_gasPrice"
}

func (p *ProxyETHGasPrice) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	revoresp, err := p.Revo.GetGasPrice(c.Request().Context())
	if err != nil {
		return nil, eth.NewCallbackError(err.Error())
	}

	// revo res -> eth res
	return p.response(revoresp), nil
}

func (p *ProxyETHGasPrice) response(revoresp *big.Int) string {
	// 34 GWEI is the minimum price that REVO will confirm tx with
	return hexutil.EncodeBig(convertFromSatoshiToWei(revoresp))
}
