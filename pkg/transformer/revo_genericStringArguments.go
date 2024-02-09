package transformer

import (
	"github.com/labstack/echo"
	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/revo"
)

type ProxyREVOGenericStringArguments struct {
	*revo.Revo
	prefix string
	method string
}

var _ ETHProxy = (*ProxyREVOGenericStringArguments)(nil)

func (p *ProxyREVOGenericStringArguments) Method() string {
	return p.prefix + "_" + p.method
}

func (p *ProxyREVOGenericStringArguments) Request(req *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	var params eth.StringsArguments
	if err := unmarshalRequest(req.Params, &params); err != nil {
		// TODO: Correct error code?
		return nil, eth.NewInvalidParamsError("couldn't unmarshal request parameters")
	}

	if len(params) != 1 {
		return nil, eth.NewInvalidParamsError("require 1 argument: the base58 Revo address")
	}

	return p.request(params)
}

func (p *ProxyREVOGenericStringArguments) request(params eth.StringsArguments) (*string, eth.JSONRPCError) {
	var response string
	err := p.Client.Request(p.method, params, &response)
	if err != nil {
		return nil, eth.NewInvalidRequestError(err.Error())
	}

	return &response, nil
}
