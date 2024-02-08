package transformer

import (
	"encoding/json"
	"testing"

	"github.com/btcsuite/btcutil"
	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/internal"
	"github.com/revolutionchain/charon/pkg/qtum"
)

func TestGetAccountInfoRequest(t *testing.T) {
	//prepare request
	requestParams := []json.RawMessage{[]byte(`"0x1e6f89d7399081b4f8f8aa1ae2805a5efff2f960"`), []byte(`"123"`)}
	requestRPC, err := internal.PrepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}
	//prepare client
	mockedClientDoer := internal.NewDoerMappedMock()
	qtumClient, err := internal.CreateMockedClient(mockedClientDoer)
	if err != nil {
		t.Fatal(err)
	}

	//prepare account
	account, err := btcutil.DecodeWIF("5JK4Gu9nxCvsCxiq9Zf3KdmA9ACza6dUn5BRLVWAYEtQabdnJ89")
	if err != nil {
		t.Fatal(err)
	}
	qtumClient.Accounts = append(qtumClient.Accounts, account)

	//prepare responses
	getAccountInfoResponse := qtum.GetAccountInfoResponse{
		Address: "1e6f89d7399081b4f8f8aa1ae2805a5efff2f960",
		Balance: 12431243,
		// Storage json.RawMessage `json:"storage"`,
		Code: "606060405236156100ad576000357c0100000000000000000...",
	}
	err = mockedClientDoer.AddResponseWithRequestID(3, qtum.MethodGetAccountInfo, getAccountInfoResponse)
	if err != nil {
		t.Fatal(err)
	}

	//preparing proxy & executing request
	proxyEth := ProxyETHGetCode{qtumClient}
	got, jsonErr := proxyEth.Request(requestRPC, internal.NewEchoContext())
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	want := eth.GetCodeResponse("0x606060405236156100ad576000357c0100000000000000000...")

	internal.CheckTestResultEthRequestRPC(*requestRPC, want, got, t, false)
}

func TestGetCodeInvalidAddressRequest(t *testing.T) {
	//prepare request
	requestParams := []json.RawMessage{[]byte(`"0x0000000000000000000000000000000000000000"`), []byte(`"123"`)}
	requestRPC, err := internal.PrepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}
	//prepare client
	mockedClientDoer := internal.NewDoerMappedMock()
	qtumClient, err := internal.CreateMockedClient(mockedClientDoer)
	if err != nil {
		t.Fatal(err)
	}

	//prepare responses
	getAccountInfoErrorResponse := qtum.GetErrorResponse(qtum.ErrInvalidAddress)
	if getAccountInfoErrorResponse == nil {
		panic("mocked error response is nil")
	}
	err = mockedClientDoer.AddError(qtum.MethodGetAccountInfo, getAccountInfoErrorResponse)
	if err != nil {
		t.Fatal(err)
	}

	//preparing proxy & executing request
	proxyEth := ProxyETHGetCode{qtumClient}
	got, jsonErr := proxyEth.Request(requestRPC, internal.NewEchoContext())
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	want := eth.GetCodeResponse("0x")

	internal.CheckTestResultEthRequestRPC(*requestRPC, want, got, t, false)
}
