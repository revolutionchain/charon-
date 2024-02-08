package transformer

import (
	"github.com/labstack/echo"
	"github.com/revolutionchain/charon/pkg/eth"
)

type ETHGetUncleByBlockHashAndIndex struct {
}

func (p *ETHGetUncleByBlockHashAndIndex) Method() string {
	return "eth_getUncleByBlockHashAndIndex"
}

func (p *ETHGetUncleByBlockHashAndIndex) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	// hardcoded to nil
	return nil, nil
}
