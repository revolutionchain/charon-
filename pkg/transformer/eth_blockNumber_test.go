package transformer

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/internal"
	"github.com/revolutionchain/charon/pkg/revo"
)

func TestBlockNumberRequest(t *testing.T) {
	//preparing request
	requestParams := []json.RawMessage{}
	request, err := internal.PrepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}

	mockedClientDoer := internal.NewDoerMappedMock()
	revoClient, err := internal.CreateMockedClient(mockedClientDoer)
	if err != nil {
		t.Fatal(err)
	}

	//preparing client response
	getBlockCountResponse := revo.GetBlockCountResponse{Int: big.NewInt(11284900)}
	err = mockedClientDoer.AddResponseWithRequestID(2, revo.MethodGetBlockCount, getBlockCountResponse)
	if err != nil {
		t.Fatal(err)
	}

	//preparing proxy & executing request
	proxyEth := ProxyETHBlockNumber{revoClient}
	got, jsonErr := proxyEth.Request(request, internal.NewEchoContext())
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	want := eth.BlockNumberResponse("0xac31a4")

	internal.CheckTestResultEthRequestRPC(*request, &want, got, t, false)
}
