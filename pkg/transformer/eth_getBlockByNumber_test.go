package transformer

import (
	"encoding/json"
	"testing"

	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/internal"
	"github.com/revolutionchain/charon/pkg/revo"
)

func initializeProxyETHGetBlockByNumber(revoClient *revo.Revo) ETHProxy {
	return &ProxyETHGetBlockByNumber{revoClient}
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
	revoClient, err := internal.CreateMockedClient(mockedClientDoer)

	unknownBlockResponse := revo.GetErrorResponse(revo.ErrInvalidParameter)
	err = mockedClientDoer.AddError(revo.MethodGetBlockHash, unknownBlockResponse)
	if err != nil {
		t.Fatal(err)
	}

	//preparing proxy & executing request
	proxyEth := ProxyETHGetBlockByNumber{revoClient}
	got, jsonErr := proxyEth.Request(request, internal.NewEchoContext())
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	want := (*eth.GetBlockByNumberResponse)(nil)

	internal.CheckTestResultDefault(want, got, t, false)
}
