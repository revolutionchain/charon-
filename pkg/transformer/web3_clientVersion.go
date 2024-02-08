package transformer

import (
	"runtime"

	"github.com/labstack/echo"
	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/params"
)

// Web3ClientVersion implements web3_clientVersion
type Web3ClientVersion struct {
	// *qtum.Qtum
}

func (p *Web3ClientVersion) Method() string {
	return "web3_clientVersion"
}

func (p *Web3ClientVersion) Request(_ *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	return "Charon/" + params.VersionWithGitSha + "/" + runtime.GOOS + "-" + runtime.GOARCH + "/" + runtime.Version(), nil
}

// func (p *Web3ClientVersion) ToResponse(ethresp *qtum.CallContractResponse) *eth.CallResponse {
// 	data := utils.AddHexPrefix(ethresp.ExecutionResult.Output)
// 	qtumresp := eth.CallResponse(data)
// 	return &qtumresp
// }
