package transformer

import (
	"github.com/labstack/echo"
	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/qtum"
)

type ProxyQTUMGenericStringArguments struct {
	*qtum.Qtum
	prefix string
	method string
}

var _ ETHProxy = (*ProxyQTUMGenericStringArguments)(nil)

func (p *ProxyQTUMGenericStringArguments) Method() string {
	return p.prefix + "_" + p.method
}

func (p *ProxyQTUMGenericStringArguments) Request(req *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	var params eth.StringsArguments
	if err := unmarshalRequest(req.Params, &params); err != nil {
		// TODO: Correct error code?
		return nil, eth.NewInvalidParamsError("couldn't unmarshal request parameters")
	}

	if len(params) != 1 {
		return nil, eth.NewInvalidParamsError("require 1 argument: the base58 Qtum address")
	}

	return p.request(params)
}

func (p *ProxyQTUMGenericStringArguments) request(params eth.StringsArguments) (*string, eth.JSONRPCError) {
	var response string
	err := p.Client.Request(p.method, params, &response)
	if err != nil {
		return nil, eth.NewInvalidRequestError(err.Error())
	}

	return &response, nil
}
