package transformer

import (
	"encoding/json"
	"testing"

	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/internal"
	"github.com/revolutionchain/charon/pkg/qtum"
)

func TestMiningRequest(t *testing.T) {
	//preparing the request
	requestParams := []json.RawMessage{} //eth_hashrate has no params
	request, err := internal.PrepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}

	mockedClientDoer := internal.NewDoerMappedMock()
	qtumClient, err := internal.CreateMockedClient(mockedClientDoer)
	if err != nil {
		t.Fatal(err)
	}

	getMiningResponse := qtum.GetMiningResponse{Staking: true}
	err = mockedClientDoer.AddResponse(qtum.MethodGetStakingInfo, getMiningResponse)
	if err != nil {
		t.Fatal(err)
	}

	proxyEth := ProxyETHMining{qtumClient}
	got, jsonErr := proxyEth.Request(request, internal.NewEchoContext())
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	want := eth.MiningResponse(true)

	internal.CheckTestResultEthRequestRPC(*request, &want, got, t, false)
}
