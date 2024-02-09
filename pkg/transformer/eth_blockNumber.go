package transformer

import (
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/revo"
)

// ProxyETHBlockNumber implements ETHProxy
type ProxyETHBlockNumber struct {
	*revo.Revo
}

func (p *ProxyETHBlockNumber) Method() string {
	return "eth_blockNumber"
}

func (p *ProxyETHBlockNumber) Request(_ *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	return p.request(c, 5)
}

func (p *ProxyETHBlockNumber) request(c echo.Context, retries int) (*eth.BlockNumberResponse, eth.JSONRPCError) {
	revoresp, err := p.Revo.GetBlockCount(c.Request().Context())
	if err != nil {
		if retries > 0 && strings.Contains(err.Error(), revo.ErrTryAgain.Error()) {
			ctx := c.Request().Context()
			t := time.NewTimer(500 * time.Millisecond)
			select {
			case <-ctx.Done():
				return nil, eth.NewCallbackError(err.Error())
			case <-t.C:
				// fallthrough
			}
			return p.request(c, retries-1)
		}
		return nil, eth.NewCallbackError(err.Error())
	}

	// revo res -> eth res
	return p.ToResponse(revoresp), nil
}

func (p *ProxyETHBlockNumber) ToResponse(revoresp *revo.GetBlockCountResponse) *eth.BlockNumberResponse {
	hexVal := hexutil.EncodeBig(revoresp.Int)
	ethresp := eth.BlockNumberResponse(hexVal)
	return &ethresp
}
