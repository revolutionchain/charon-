package transformer

import (
	"encoding/json"
	"testing"

	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/internal"
	"github.com/revolutionchain/charon/pkg/revo"
)

func TestMiningRequest(t *testing.T) {
	//preparing the request
	requestParams := []json.RawMessage{} //eth_hashrate has no params
	request, err := internal.PrepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}

	mockedClientDoer := internal.NewDoerMappedMock()
	revoClient, err := internal.CreateMockedClient(mockedClientDoer)
	if err != nil {
		t.Fatal(err)
	}

	getMiningResponse := revo.GetMiningResponse{Staking: true}
	err = mockedClientDoer.AddResponse(revo.MethodGetStakingInfo, getMiningResponse)
	if err != nil {
		t.Fatal(err)
	}

	proxyEth := ProxyETHMining{revoClient}
	got, jsonErr := proxyEth.Request(request, internal.NewEchoContext())
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	want := eth.MiningResponse(true)

	internal.CheckTestResultEthRequestRPC(*request, &want, got, t, false)
}
