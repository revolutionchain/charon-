package transformer

import (
	"context"

	"github.com/labstack/echo"
	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/revo"
	"github.com/revolutionchain/charon/pkg/utils"
)

// ProxyETHGetCode implements ETHProxy
type ProxyETHGetCode struct {
	*revo.Revo
}

func (p *ProxyETHGetCode) Method() string {
	return "eth_getCode"
}

func (p *ProxyETHGetCode) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	var req eth.GetCodeRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		// TODO: Correct error code?
		return nil, eth.NewInvalidParamsError(err.Error())
	}

	return p.request(c.Request().Context(), &req)
}

func (p *ProxyETHGetCode) request(ctx context.Context, ethreq *eth.GetCodeRequest) (eth.GetCodeResponse, eth.JSONRPCError) {
	revoreq := revo.GetAccountInfoRequest(utils.RemoveHexPrefix(ethreq.Address))

	revoresp, err := p.GetAccountInfo(ctx, &revoreq)
	if err != nil {
		if err == revo.ErrInvalidAddress {
			/**
			// correct response for an invalid address
			{
				"jsonrpc": "2.0",
				"id": 123,
				"result": "0x"
			}
			**/
			return "0x", nil
		} else {
			return "", eth.NewCallbackError(err.Error())
		}
	}

	// revo res -> eth res
	return eth.GetCodeResponse(utils.AddHexPrefix(revoresp.Code)), nil
}
