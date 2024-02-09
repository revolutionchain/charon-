package transformer

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/revo"
)

// ProxyETHUninstallFilter implements ETHProxy
type ProxyETHUninstallFilter struct {
	*revo.Revo
	filter *eth.FilterSimulator
}

func (p *ProxyETHUninstallFilter) Method() string {
	return "eth_uninstallFilter"
}

func (p *ProxyETHUninstallFilter) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	var req eth.UninstallFilterRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		// TODO: Correct error code?
		return nil, eth.NewInvalidParamsError(err.Error())
	}

	return p.request(&req)
}

func (p *ProxyETHUninstallFilter) request(ethreq *eth.UninstallFilterRequest) (eth.UninstallFilterResponse, eth.JSONRPCError) {
	id, err := hexutil.DecodeUint64(string(*ethreq))
	if err != nil {
		return false, eth.NewInvalidParamsError(err.Error())
	}

	// uninstall
	p.filter.Uninstall(id)

	return true, nil
}
