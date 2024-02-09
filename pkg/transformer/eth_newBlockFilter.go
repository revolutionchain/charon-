package transformer

import (
	"context"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/revo"
)

// ProxyETHNewBlockFilter implements ETHProxy
type ProxyETHNewBlockFilter struct {
	*revo.Revo
	filter *eth.FilterSimulator
}

func (p *ProxyETHNewBlockFilter) Method() string {
	return "eth_newBlockFilter"
}

func (p *ProxyETHNewBlockFilter) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	return p.request(c.Request().Context())
}

func (p *ProxyETHNewBlockFilter) request(ctx context.Context) (eth.NewBlockFilterResponse, eth.JSONRPCError) {
	blockCount, err := p.GetBlockCount(ctx)
	if err != nil {
		return "", eth.NewCallbackError(err.Error())
	}

	filter := p.filter.New(eth.NewBlockFilterTy)
	filter.Data.Store("lastBlockNumber", blockCount.Uint64())

	p.GenerateIfPossible()

	return eth.NewBlockFilterResponse(hexutil.EncodeUint64(filter.ID)), nil
}
