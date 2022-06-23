package transformer

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/qtum"
)

type ProxyETHChainId struct {
	*qtum.Qtum
}

func (p *ProxyETHChainId) Method() string {
	return "eth_chainId"
}

func (p *ProxyETHChainId) Request(req *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	chainId, err := getChainId(p.Qtum)
	if err != nil {
		return nil, err
	}
	return eth.ChainIdResponse(hexutil.EncodeBig(chainId)), nil
}

func getChainId(p *qtum.Qtum) (*big.Int, eth.JSONRPCError) {
	var qtumresp *qtum.GetBlockChainInfoResponse
	if err := p.Request(qtum.MethodGetBlockChainInfo, nil, &qtumresp); err != nil {
		return nil, eth.NewCallbackError(err.Error())
	}

	var chainId *big.Int
	switch strings.ToLower(qtumresp.Chain) {
	case "main":
		chainId = big.NewInt(81)
	case "test":
		chainId = big.NewInt(8889)
	case "regtest":
		chainId = big.NewInt(8890)
	default:
		chainId = big.NewInt(8890)
		p.GetDebugLogger().Log("msg", "Unknown chain "+qtumresp.Chain)
	}

	return chainId, nil
}
