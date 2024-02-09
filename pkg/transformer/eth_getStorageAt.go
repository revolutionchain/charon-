package transformer

import (
	"context"
	"fmt"

	"github.com/labstack/echo"
	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/revo"
	"github.com/revolutionchain/charon/pkg/utils"
)

// ProxyETHGetStorageAt implements ETHProxy
type ProxyETHGetStorageAt struct {
	*revo.Revo
}

func (p *ProxyETHGetStorageAt) Method() string {
	return "eth_getStorageAt"
}

func (p *ProxyETHGetStorageAt) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	var req eth.GetStorageRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		// TODO: Correct error code?
		return nil, eth.NewInvalidParamsError(err.Error())
	}

	revoAddress := utils.RemoveHexPrefix(req.Address)
	blockNumber, err := getBlockNumberByParam(c.Request().Context(), p.Revo, req.BlockNumber, false)
	if err != nil {
		p.GetDebugLogger().Log("msg", fmt.Sprintf("Failed to get block number by param for '%s'", req.BlockNumber), "err", err)
		return nil, err
	}

	return p.request(
		c.Request().Context(),
		&revo.GetStorageRequest{
			Address:     revoAddress,
			BlockNumber: blockNumber,
		},
		utils.RemoveHexPrefix(req.Index),
	)
}

func (p *ProxyETHGetStorageAt) request(ctx context.Context, ethreq *revo.GetStorageRequest, index string) (*eth.GetStorageResponse, eth.JSONRPCError) {
	revoresp, err := p.Revo.GetStorage(ctx, ethreq)
	if err != nil {
		return nil, eth.NewCallbackError(err.Error())
	}

	// revo res -> eth res
	return p.ToResponse(revoresp, index), nil
}

func (p *ProxyETHGetStorageAt) ToResponse(revoresp *revo.GetStorageResponse, slot string) *eth.GetStorageResponse {
	// the value for unknown anything
	storageData := eth.GetStorageResponse("0x0000000000000000000000000000000000000000000000000000000000000000")
	if len(slot) != 64 {
		slot = leftPadStringWithZerosTo64Bytes(slot)
	}
	for _, outerValue := range *revoresp {
		revoStorageData, ok := outerValue[slot]
		if ok {
			storageData = eth.GetStorageResponse(utils.AddHexPrefix(revoStorageData))
			return &storageData
		}
	}

	return &storageData
}

// left pad a string with leading zeros to fit 64 bytes
func leftPadStringWithZerosTo64Bytes(hex string) string {
	return fmt.Sprintf("%064v", hex)
}
