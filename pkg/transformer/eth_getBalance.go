package transformer

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/labstack/echo"
	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/revo"
	"github.com/revolutionchain/charon/pkg/utils"
)

// ProxyETHGetBalance implements ETHProxy
type ProxyETHGetBalance struct {
	*revo.Revo
}

func (p *ProxyETHGetBalance) Method() string {
	return "eth_getBalance"
}

func (p *ProxyETHGetBalance) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	var req eth.GetBalanceRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		// TODO: Correct error code?
		return nil, eth.NewInvalidParamsError(err.Error())
	}

	addr := utils.RemoveHexPrefix(req.Address)
	{
		// is address a contract or an account?
		revoreq := revo.GetAccountInfoRequest(addr)
		revoresp, err := p.GetAccountInfo(c.Request().Context(), &revoreq)

		// the address is a contract
		if err == nil {
			// the unit of the balance Satoshi
			p.GetDebugLogger().Log("method", p.Method(), "address", req.Address, "msg", "is a contract")
			return hexutil.EncodeUint64(uint64(revoresp.Balance)), nil
		}
	}

	{
		// try account
		base58Addr, err := p.FromHexAddress(addr)
		if err != nil {
			p.GetDebugLogger().Log("method", p.Method(), "address", req.Address, "msg", "error parsing address", "error", err)
			return nil, eth.NewCallbackError(err.Error())
		}

		revoreq := revo.GetAddressBalanceRequest{Address: base58Addr}
		revoresp, err := p.GetAddressBalance(c.Request().Context(), &revoreq)
		if err != nil {
			if err == revo.ErrInvalidAddress {
				// invalid address should return 0x0
				return "0x0", nil
			}
			p.GetDebugLogger().Log("method", p.Method(), "address", req.Address, "msg", "error getting address balance", "error", err)
			return nil, eth.NewCallbackError(err.Error())
		}

		// 1 REVO = 10 ^ 8 Satoshi
		balance := new(big.Int).SetUint64(revoresp.Balance)

		//Balance for ETH response is represented in Weis (1 REVO Satoshi = 10 ^ 10 Wei)
		balance = balance.Mul(balance, big.NewInt(10000000000))

		return hexutil.EncodeBig(balance), nil
	}
}
