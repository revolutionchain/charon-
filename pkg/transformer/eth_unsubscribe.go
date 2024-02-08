package transformer

import (
	"github.com/labstack/echo"
	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/notifier"
	"github.com/revolutionchain/charon/pkg/qtum"
)

// ETHUnsubscribe implements ETHProxy
type ETHUnsubscribe struct {
	*qtum.Qtum
	*notifier.Agent
}

func (p *ETHUnsubscribe) Method() string {
	return "eth_unsubscribe"
}

func (p *ETHUnsubscribe) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	notifier := getNotifier(c)
	if notifier == nil {
		p.GetLogger().Log("msg", "eth_unsubscribe only supported over websocket")
		/*
			// TODO
			{
				"jsonrpc": "2.0",
				"id": 580,
				"error": {
					"code": -32601,
					"message": "The method eth_unsubscribe does not exist/is not available"
				}
			}
		*/
		return nil, eth.NewMethodNotFoundError("eth_subscribe")
	}

	var req eth.EthUnsubscribeRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		// TODO: Correct error code?
		return nil, eth.NewInvalidParamsError(err.Error())
	}

	return p.request(&req, notifier)
}

func (p *ETHUnsubscribe) request(req *eth.EthUnsubscribeRequest, notifier *notifier.Notifier) (eth.EthUnsubscribeResponse, eth.JSONRPCError) {
	if len(*req) != 1 {
		// TODO: Correct error code?
		return false, eth.NewInvalidParamsError("requires one parameter")
	}
	param := (*req)[0]
	success := notifier.Unsubscribe(param)
	return eth.EthUnsubscribeResponse(success), nil
}
