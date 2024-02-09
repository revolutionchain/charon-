package transformer

import (
	"context"

	"github.com/labstack/echo"
	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/revo"
)

// ProxyETHGetHashrate implements ETHProxy
type ProxyETHMining struct {
	*revo.Revo
}

func (p *ProxyETHMining) Method() string {
	return "eth_mining"
}

func (p *ProxyETHMining) Request(_ *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	return p.request(c.Request().Context())
}

func (p *ProxyETHMining) request(ctx context.Context) (*eth.MiningResponse, eth.JSONRPCError) {
	revoresp, err := p.Revo.GetMining(ctx)
	if err != nil {
		return nil, eth.NewCallbackError(err.Error())
	}

	// revo res -> eth res
	return p.ToResponse(revoresp), nil
}

func (p *ProxyETHMining) ToResponse(revoresp *revo.GetMiningResponse) *eth.MiningResponse {
	ethresp := eth.MiningResponse(revoresp.Staking)
	return &ethresp
}
