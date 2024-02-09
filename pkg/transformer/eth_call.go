package transformer

import (
	"context"
	"math/big"

	"github.com/labstack/echo"
	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/revo"
	"github.com/revolutionchain/charon/pkg/utils"
)

// ProxyETHCall implements ETHProxy
type ProxyETHCall struct {
	*revo.Revo
}

func (p *ProxyETHCall) Method() string {
	return "eth_call"
}

func (p *ProxyETHCall) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	var req eth.CallRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		// TODO: Is this correct error code?
		return nil, eth.NewInvalidParamsError(err.Error())
	}

	return p.request(c.Request().Context(), &req)
}

func (p *ProxyETHCall) request(ctx context.Context, ethreq *eth.CallRequest) (interface{}, eth.JSONRPCError) {
	// eth req -> revo req
	revoreq, jsonErr := p.ToRequest(ethreq)
	if jsonErr != nil {
		return nil, jsonErr
	}
	if revoreq.GasLimit != nil && revoreq.GasLimit.Cmp(big.NewInt(40000000)) > 0 {
		revoresp := eth.CallResponse("0x")
		p.Revo.GetLogger().Log("msg", "Caller gas above allowance, capping", "requested", revoreq.GasLimit.Int64(), "cap", "40,000,000")
		return &revoresp, nil
	}

	revoresp, err := p.CallContract(ctx, revoreq)
	if err != nil {
		if err == revo.ErrInvalidAddress {
			revoresp := eth.CallResponse("0x")
			return &revoresp, nil
		}

		return nil, eth.NewCallbackError(err.Error())
	}

	// revo res -> eth res
	return p.ToResponse(revoresp), nil
}

func (p *ProxyETHCall) ToRequest(ethreq *eth.CallRequest) (*revo.CallContractRequest, eth.JSONRPCError) {
	from := ethreq.From
	var err error
	if utils.IsEthHexAddress(from) {
		from, err = p.FromHexAddress(from)
		if err != nil {
			return nil, eth.NewCallbackError(err.Error())
		}
	}

	var gasLimit *big.Int
	if ethreq.Gas != nil {
		gasLimit = ethreq.Gas.Int
	}

	if gasLimit != nil && gasLimit.Int64() < MinimumGasLimit {
		p.GetLogger().Log("msg", "Gas limit is too low", "gasLimit", gasLimit.String())
	}

	return &revo.CallContractRequest{
		To:       ethreq.To,
		From:     from,
		Data:     ethreq.Data,
		GasLimit: gasLimit,
	}, nil
}

func (p *ProxyETHCall) ToResponse(qresp *revo.CallContractResponse) interface{} {
	if qresp.ExecutionResult.Output == "" {
		return eth.NewJSONRPCError(
			-32000,
			"Revert: executionResult output is empty",
			nil,
		)
	}

	data := utils.AddHexPrefix(qresp.ExecutionResult.Output)
	revoresp := eth.CallResponse(data)
	return &revoresp

}
