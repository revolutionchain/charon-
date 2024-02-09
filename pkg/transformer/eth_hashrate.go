package transformer

import (
	"context"
	"math"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/revo"
)

// ProxyETHGetHashrate implements ETHProxy
type ProxyETHHashrate struct {
	*revo.Revo
}

func (p *ProxyETHHashrate) Method() string {
	return "eth_hashrate"
}

func (p *ProxyETHHashrate) Request(_ *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	return p.request(c.Request().Context())
}

func (p *ProxyETHHashrate) request(ctx context.Context) (*eth.HashrateResponse, eth.JSONRPCError) {
	revoresp, err := p.Revo.GetHashrate(ctx)
	if err != nil {
		return nil, eth.NewCallbackError(err.Error())
	}

	// revo res -> eth res
	return p.ToResponse(revoresp), nil
}

func (p *ProxyETHHashrate) ToResponse(revoresp *revo.GetHashrateResponse) *eth.HashrateResponse {
	hexVal := hexutil.EncodeUint64(math.Float64bits(revoresp.Difficulty))
	ethresp := eth.HashrateResponse(hexVal)
	return &ethresp
}
