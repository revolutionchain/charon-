package transformer

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/revo"
)

type ProxyETHChainId struct {
	*revo.Revo
}

func (p *ProxyETHChainId) Method() string {
	return "eth_chainId"
}

func (p *ProxyETHChainId) Request(req *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	chainId, err := getChainId(p.Revo)
	if err != nil {
		return nil, err
	}
	return eth.ChainIdResponse(hexutil.EncodeBig(chainId)), nil
}

func getChainId(p *revo.Revo) (*big.Int, eth.JSONRPCError) {
	return big.NewInt(int64(p.ChainId())), nil
}
