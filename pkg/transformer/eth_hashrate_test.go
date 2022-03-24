package transformer

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/internal"
	"github.com/qtumproject/janus/pkg/qtum"
)

func TestHashrateRequest(t *testing.T) {
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

	exampleResponse := `{"enabled": true, "staking": false, "errors": "", "currentblocktx": 0, "pooledtx": 0, "difficulty": 4.656542373906925e-010, "search-interval": 0, "weight": 0, "netstakeweight": 0, "expectedtime": 0}`
	getHashrateResponse := qtum.GetHashrateResponse{}
	unmarshalRequest([]byte(exampleResponse), &getHashrateResponse)

	err = mockedClientDoer.AddResponse(qtum.MethodGetStakingInfo, getHashrateResponse)
	if err != nil {
		t.Fatal(err)
	}

	proxyEth := ProxyETHHashrate{qtumClient}
	got, jsonErr := proxyEth.Request(request, nil)
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	expected := hexutil.EncodeUint64(math.Float64bits(4.656542373906925e-010))
	want := eth.HashrateResponse(expected)

	internal.CheckTestResultEthRequestRPC(*request, &want, got, t, false)
}
