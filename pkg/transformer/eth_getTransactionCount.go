package transformer

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/qtum"
)

// ProxyETHEstimateGas implements ETHProxy
type ProxyETHTxCount struct {
	*qtum.Qtum
}

func (p *ProxyETHTxCount) Method() string {
	return "eth_getTransactionCount"
}

func (p *ProxyETHTxCount) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {

	/* not sure we need this. Need to figure out how to best unmarshal this in the future. For now this will work.
	var req eth.GetTransactionCountRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		return nil, err
	}*/
	qtumresp, err := p.Qtum.GetTransactionCount(c.Request().Context(), "", "")
	if err != nil {
		return nil, eth.NewCallbackError(err.Error())
	}

	// qtum res -> eth res
	return p.response(qtumresp), nil
}

func (p *ProxyETHTxCount) response(qtumresp *big.Int) string {
	return hexutil.EncodeBig(qtumresp)
}
