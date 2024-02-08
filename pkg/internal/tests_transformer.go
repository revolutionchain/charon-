package internal

import (
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/revolutionchain/charon/pkg/eth"
)

type ETHProxy interface {
	Request(*eth.JSONRPCRequest, echo.Context) (interface{}, eth.JSONRPCError)
	Method() string
}

type mockTransformer struct {
	proxies map[string]ETHProxy
}

func (t *mockTransformer) Transform(req *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	proxy, ok := t.proxies[req.Method]
	if !ok {
		return nil, eth.NewCallbackError("couldn't get proxy")
	}
	resp, err := proxy.Request(req, c)
	if err != nil {
		return nil, eth.NewCallbackError(errors.WithMessagef(err.Error(), "couldn't proxy %s request", req.Method).Error())
	}
	return resp, nil
}

func newTransformer(proxies []ETHProxy) *mockTransformer {
	t := &mockTransformer{
		proxies: make(map[string]ETHProxy),
	}

	for _, proxy := range proxies {
		t.proxies[proxy.Method()] = proxy
	}

	return t
}

func NewMockTransformer(proxies []ETHProxy) *mockTransformer {
	return newTransformer(proxies)
}

type mockETHProxy struct {
	method   string
	response interface{}
}

func NewMockETHProxy(method string, response interface{}) ETHProxy {
	return &mockETHProxy{
		method:   method,
		response: response,
	}
}

func (e *mockETHProxy) Request(*eth.JSONRPCRequest, echo.Context) (interface{}, eth.JSONRPCError) {
	return e.response, nil
}

func (e *mockETHProxy) Method() string {
	return e.method
}
