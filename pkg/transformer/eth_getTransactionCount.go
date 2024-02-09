package transformer

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/revo"
)

// ProxyETHEstimateGas implements ETHProxy
type ProxyETHTxCount struct {
	*revo.Revo
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
	revoresp, err := p.Revo.GetTransactionCount(c.Request().Context(), "", "")
	if err != nil {
		return nil, eth.NewCallbackError(err.Error())
	}

	// revo res -> eth res
	return p.response(revoresp), nil
}

func (p *ProxyETHTxCount) response(revoresp *big.Int) string {
	return hexutil.EncodeBig(revoresp)
}
