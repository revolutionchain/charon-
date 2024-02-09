package transformer

import (
	"reflect"
	"strconv"

	"github.com/labstack/echo"
	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/revo"
	"github.com/revolutionchain/charon/pkg/utils"
)

type ProxyREVOGenerateToAddress struct {
	*revo.Revo
}

var _ ETHProxy = (*ProxyREVOGenerateToAddress)(nil)

func (p *ProxyREVOGenerateToAddress) Method() string {
	return "dev_generatetoaddress"
}

func (p *ProxyREVOGenerateToAddress) Request(req *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	if !p.CanGenerate() {
		return nil, eth.NewInvalidRequestError("Can only generate on regtest")
	}

	var params []interface{}
	if err := unmarshalRequest(req.Params, &params); err != nil {
		// TODO: Correct error code?
		return nil, eth.NewInvalidParamsError("couldn't unmarshal request parameters")
	}

	if len(params) != 2 {
		return nil, eth.NewInvalidParamsError("require 2 arguments: blocks, the base58/hex address to mine rewards to")
	}

	return p.request(params)
}

func (p *ProxyREVOGenerateToAddress) request(params []interface{}) (*[]string, eth.JSONRPCError) {
	blocks := params[0]
	generateTo, ok := params[1].(string)
	if !ok {
		return nil, eth.NewInvalidParamsError("second paramter must be string")
	}

	var blocksInteger int64
	var err error

	if blocksString, ok := blocks.(string); ok {
		blocksInteger, err = strconv.ParseInt(blocksString, 10, 64)
		if err != nil {
			return nil, eth.NewInvalidParamsError("Couldn't parse blocks")
		}
	} else if blocksNumber, ok := blocks.(float64); ok {
		blocksInteger = int64(blocksNumber)
	} else {
		return nil, eth.NewInvalidParamsError("Unknown blocks type: " + reflect.TypeOf(blocks).String())
	}

	if blocksInteger <= 0 {
		return nil, eth.NewInvalidParamsError("Blocks to generate must be > 0")
	}

	hex := utils.RemoveHexPrefix(generateTo)
	base58Address, err := p.FromHexAddress(hex)
	if err != nil {
		// already base58?
		base58Address = generateTo
	}

	var response []string
	err = p.Client.Request(revo.MethodGenerateToAddress, []interface{}{blocksInteger, base58Address}, &response)
	if err != nil {
		return nil, eth.NewInvalidRequestError(err.Error())
	}

	return &response, nil
}
