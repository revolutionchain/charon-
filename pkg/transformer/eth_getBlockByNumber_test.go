package transformer

import (
	"encoding/json"
	"testing"

	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/internal"
	"github.com/qtumproject/janus/pkg/qtum"
)

func initializeProxyETHGetBlockByNumber(qtumClient *qtum.Qtum) ETHProxy {
	return &ProxyETHGetBlockByNumber{qtumClient}
}

func TestGetBlockByNumberRequest(t *testing.T) {
	testETHProxyRequest(
		t,
		initializeProxyETHGetBlockByNumber,
		[]json.RawMessage{[]byte(`"` + internal.GetTransactionByHashBlockNumberHex + `"`), []byte(`false`)},
		&internal.GetTransactionByHashResponse,
	)
}

func TestGetBlockByNumberWithTransactionsRequest(t *testing.T) {
	testETHProxyRequest(
		t,
		initializeProxyETHGetBlockByNumber,
		[]json.RawMessage{[]byte(`"` + internal.GetTransactionByHashBlockNumberHex + `"`), []byte(`true`)},
		&internal.GetTransactionByHashResponseWithTransactions,
	)
}

func TestGetBlockByNumberUnknownBlockRequest(t *testing.T) {
	requestParams := []json.RawMessage{[]byte(`"` + internal.GetTransactionByHashBlockNumberHex + `"`), []byte(`true`)}
	request, err := internal.PrepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}

	mockedClientDoer := internal.NewDoerMappedMock()
	qtumClient, err := internal.CreateMockedClient(mockedClientDoer)

	unknownBlockResponse := qtum.GetErrorResponse(qtum.ErrInvalidParameter)
	err = mockedClientDoer.AddError(qtum.MethodGetBlockHash, unknownBlockResponse)
	if err != nil {
		t.Fatal(err)
	}

	//preparing proxy & executing request
	proxyEth := ProxyETHGetBlockByNumber{qtumClient}
	got, jsonErr := proxyEth.Request(request, nil)
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	want := (*eth.GetBlockByNumberResponse)(nil)

	internal.CheckTestResultDefault(want, got, t, false)
}
