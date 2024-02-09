package transformer

import (
	"context"

	"github.com/labstack/echo"
	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/revo"
	"github.com/revolutionchain/charon/pkg/utils"
)

// ProxyETHSendRawTransaction implements ETHProxy
type ProxyETHSendRawTransaction struct {
	*revo.Revo
}

var _ ETHProxy = (*ProxyETHSendRawTransaction)(nil)

func (p *ProxyETHSendRawTransaction) Method() string {
	return "eth_sendRawTransaction"
}

func (p *ProxyETHSendRawTransaction) Request(req *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	var params eth.SendRawTransactionRequest
	if err := unmarshalRequest(req.Params, &params); err != nil {
		// TODO: Correct error code?
		return nil, eth.NewInvalidParamsError(err.Error())
	}
	if params[0] == "" {
		// TODO: Correct error code?
		return nil, eth.NewInvalidParamsError("invalid parameter: raw transaction hexed string is empty")
	}

	return p.request(c.Request().Context(), params)
}

func (p *ProxyETHSendRawTransaction) request(ctx context.Context, params eth.SendRawTransactionRequest) (eth.SendRawTransactionResponse, eth.JSONRPCError) {
	var (
		revoHexedRawTx = utils.RemoveHexPrefix(params[0])
		req            = revo.SendRawTransactionRequest([1]string{revoHexedRawTx})
	)

	revoresp, err := p.Revo.SendRawTransaction(ctx, &req)
	if err != nil {
		if err == revo.ErrVerifyAlreadyInChain {
			// already committed
			// we need to send back the tx hash
			rawTx, err := p.Revo.DecodeRawTransaction(ctx, revoHexedRawTx)
			if err != nil {
				p.GetErrorLogger().Log("msg", "Error decoding raw transaction for duplicate raw transaction", "err", err)
				return eth.SendRawTransactionResponse(""), eth.NewCallbackError(err.Error())
			}
			revoresp = &revo.SendRawTransactionResponse{Result: rawTx.Hash}
		} else {
			return eth.SendRawTransactionResponse(""), eth.NewCallbackError(err.Error())
		}
	} else {
		p.GenerateIfPossible()
	}

	resp := *revoresp
	ethHexedTxHash := utils.AddHexPrefix(resp.Result)
	return eth.SendRawTransactionResponse(ethHexedTxHash), nil
}
