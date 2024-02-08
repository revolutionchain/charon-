package transformer

import (
	"github.com/labstack/echo"
	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/notifier"
	"github.com/revolutionchain/charon/pkg/qtum"
)

// ETHSubscribe implements ETHProxy
type ETHSubscribe struct {
	*qtum.Qtum
	*notifier.Agent
}

func (p *ETHSubscribe) Method() string {
	return "eth_subscribe"
}

func (p *ETHSubscribe) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	notifier := getNotifier(c)
	if notifier == nil {
		p.GetLogger().Log("msg", "eth_subscribe only supported over websocket")
		/*
			// TODO
			{
				"jsonrpc": "2.0",
				"id": 580,
				"error": {
					"code": -32601,
					"message": "The method eth_subscribe does not exist/is not available"
				}
			}
		*/
		return nil, eth.NewMethodNotFoundError("eth_subscribe")
	}

	var req eth.EthSubscriptionRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		// TODO: Correct error code?
		return nil, eth.NewInvalidParamsError(err.Error())
	}

	return p.request(&req, notifier)
}

func (p *ETHSubscribe) request(req *eth.EthSubscriptionRequest, notifier *notifier.Notifier) (*eth.EthSubscriptionResponse, eth.JSONRPCError) {
	notifier.ResponseRequired()
	id, err := p.NewSubscription(notifier, req)
	response := eth.EthSubscriptionResponse(id)
	if err == nil {
		return &response, nil
	} else {
		return &response, eth.NewCallbackError(err.Error())
	}
}
