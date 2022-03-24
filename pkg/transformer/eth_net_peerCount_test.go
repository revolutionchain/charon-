package transformer

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/internal"
	"github.com/qtumproject/janus/pkg/qtum"
)

func TestPeerCountRequest(t *testing.T) {
	for i := 0; i < 10; i++ {
		testDesc := fmt.Sprintf("#%d", i)
		t.Run(testDesc, func(t *testing.T) {
			testPeerCountRequest(t, i)
		})
	}
}

func testPeerCountRequest(t *testing.T, clients int) {
	//preparing the request
	requestParams := []json.RawMessage{} //net_peerCount has no params
	request, err := internal.PrepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}

	mockedClientDoer := internal.NewDoerMappedMock()
	qtumClient, err := internal.CreateMockedClient(mockedClientDoer)
	if err != nil {
		t.Fatal(err)
	}

	getPeerInfoResponse := []qtum.GetPeerInfoResponse{}
	for i := 0; i < clients; i++ {
		getPeerInfoResponse = append(getPeerInfoResponse, qtum.GetPeerInfoResponse{})
	}
	err = mockedClientDoer.AddResponseWithRequestID(2, qtum.MethodGetPeerInfo, getPeerInfoResponse)
	if err != nil {
		t.Fatal(err)
	}

	proxyEth := ProxyNetPeerCount{qtumClient}
	got, jsonErr := proxyEth.Request(request, nil)
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	want := eth.NetPeerCountResponse(hexutil.EncodeUint64(uint64(clients)))

	internal.CheckTestResultUnspecifiedInput(fmt.Sprint(clients), &want, got, t, false)
}
