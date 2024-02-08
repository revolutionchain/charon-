package transformer

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/qtum"
)

// ProxyETHEstimateGas implements ETHProxy
type ProxyETHGasPrice struct {
	*qtum.Qtum
}

func (p *ProxyETHGasPrice) Method() string {
	return "eth_gasPrice"
}

func (p *ProxyETHGasPrice) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	qtumresp, err := p.Qtum.GetGasPrice(c.Request().Context())
	if err != nil {
		return nil, eth.NewCallbackError(err.Error())
	}

	// qtum res -> eth res
	return p.response(qtumresp), nil
}

func (p *ProxyETHGasPrice) response(qtumresp *big.Int) string {
	// 34 GWEI is the minimum price that QTUM will confirm tx with
	return hexutil.EncodeBig(convertFromSatoshiToWei(qtumresp))
}
