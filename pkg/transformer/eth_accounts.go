package transformer

import (
	"github.com/labstack/echo"
	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/revo"
	"github.com/revolutionchain/charon/pkg/utils"
)

// ProxyETHAccounts implements ETHProxy
type ProxyETHAccounts struct {
	*revo.Revo
}

func (p *ProxyETHAccounts) Method() string {
	return "eth_accounts"
}

func (p *ProxyETHAccounts) Request(_ *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	return p.request()
}

func (p *ProxyETHAccounts) request() (eth.AccountsResponse, eth.JSONRPCError) {
	var accounts eth.AccountsResponse

	for _, acc := range p.Accounts {
		acc := revo.Account{acc}
		addr := acc.ToHexAddress()

		accounts = append(accounts, utils.AddHexPrefix(addr))
	}

	return accounts, nil
}

func (p *ProxyETHAccounts) ToResponse(ethresp *revo.CallContractResponse) *eth.CallResponse {
	data := utils.AddHexPrefix(ethresp.ExecutionResult.Output)
	revoresp := eth.CallResponse(data)
	return &revoresp
}
