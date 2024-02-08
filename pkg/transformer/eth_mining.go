package transformer

import (
	"context"

	"github.com/labstack/echo"
	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/qtum"
)

// ProxyETHGetHashrate implements ETHProxy
type ProxyETHMining struct {
	*qtum.Qtum
}

func (p *ProxyETHMining) Method() string {
	return "eth_mining"
}

func (p *ProxyETHMining) Request(_ *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	return p.request(c.Request().Context())
}

func (p *ProxyETHMining) request(ctx context.Context) (*eth.MiningResponse, eth.JSONRPCError) {
	qtumresp, err := p.Qtum.GetMining(ctx)
	if err != nil {
		return nil, eth.NewCallbackError(err.Error())
	}

	// qtum res -> eth res
	return p.ToResponse(qtumresp), nil
}

func (p *ProxyETHMining) ToResponse(qtumresp *qtum.GetMiningResponse) *eth.MiningResponse {
	ethresp := eth.MiningResponse(qtumresp.Staking)
	return &ethresp
}
