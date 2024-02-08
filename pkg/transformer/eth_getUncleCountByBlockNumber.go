package transformer

import (
	"github.com/labstack/echo"
	"github.com/revolutionchain/charon/pkg/eth"
)

type ETHGetUncleCountByBlockNumber struct {
}

func (p *ETHGetUncleCountByBlockNumber) Method() string {
	return "eth_getUncleCountByBlockNumber"
}

func (p *ETHGetUncleCountByBlockNumber) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	// hardcoded to 0
	return "0x0", nil
}
