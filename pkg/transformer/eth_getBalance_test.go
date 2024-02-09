package transformer

import (
	"encoding/json"
	"testing"

	"github.com/btcsuite/btcutil"
	"github.com/revolutionchain/charon/pkg/internal"
	"github.com/revolutionchain/charon/pkg/revo"
)

func TestGetBalanceRequestAccount(t *testing.T) {
	//prepare request
	requestParams := []json.RawMessage{[]byte(`"0x1e6f89d7399081b4f8f8aa1ae2805a5efff2f960"`), []byte(`"123"`)}
	requestRPC, err := internal.PrepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}
	//prepare client
	mockedClientDoer := internal.NewDoerMappedMock()
	revoClient, err := internal.CreateMockedClient(mockedClientDoer)
	if err != nil {
		t.Fatal(err)
	}

	//prepare account
	account, err := btcutil.DecodeWIF("5JK4Gu9nxCvsCxiq9Zf3KdmA9ACza6dUn5BRLVWAYEtQabdnJ89")
	if err != nil {
		t.Fatal(err)
	}
	revoClient.Accounts = append(revoClient.Accounts, account)

	//prepare responses
	fromHexAddressResponse := revo.FromHexAddressResponse("5JK4Gu9nxCvsCxiq9Zf3KdmA9ACza6dUn5BRLVWAYEtQabdnJ89")
	err = mockedClientDoer.AddResponseWithRequestID(2, revo.MethodFromHexAddress, fromHexAddressResponse)
	if err != nil {
		t.Fatal(err)
	}

	getAddressBalanceResponse := revo.GetAddressBalanceResponse{Balance: uint64(100000000), Received: uint64(100000000), Immature: int64(0)}
	err = mockedClientDoer.AddResponseWithRequestID(3, revo.MethodGetAddressBalance, getAddressBalanceResponse)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: Need getaccountinfo to return an account for unit test
	// if getaccountinfo returns nil
	// then address is contract, else account

	//preparing proxy & executing request
	proxyEth := ProxyETHGetBalance{revoClient}
	got, jsonErr := proxyEth.Request(requestRPC, internal.NewEchoContext())
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	want := string("0xde0b6b3a7640000") //1 Revo represented in Wei

	internal.CheckTestResultEthRequestRPC(*requestRPC, want, got, t, false)
}

func TestGetBalanceRequestContract(t *testing.T) {
	//prepare request
	requestParams := []json.RawMessage{[]byte(`"0x1e6f89d7399081b4f8f8aa1ae2805a5efff2f960"`), []byte(`"123"`)}
	requestRPC, err := internal.PrepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}
	//prepare client
	mockedClientDoer := internal.NewDoerMappedMock()
	revoClient, err := internal.CreateMockedClient(mockedClientDoer)
	if err != nil {
		t.Fatal(err)
	}

	//prepare account
	account, err := btcutil.DecodeWIF("5JK4Gu9nxCvsCxiq9Zf3KdmA9ACza6dUn5BRLVWAYEtQabdnJ89")
	if err != nil {
		t.Fatal(err)
	}
	revoClient.Accounts = append(revoClient.Accounts, account)

	//prepare responses
	getAccountInfoResponse := revo.GetAccountInfoResponse{
		Address: "1e6f89d7399081b4f8f8aa1ae2805a5efff2f960",
		Balance: 12431243,
		// Storage json.RawMessage `json:"storage"`,
		// Code    string          `json:"code"`,
	}
	err = mockedClientDoer.AddResponseWithRequestID(3, revo.MethodGetAccountInfo, getAccountInfoResponse)
	if err != nil {
		t.Fatal(err)
	}

	//preparing proxy & executing request
	proxyEth := ProxyETHGetBalance{revoClient}
	got, jsonErr := proxyEth.Request(requestRPC, internal.NewEchoContext())
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	want := string("0xbdaf8b")

	internal.CheckTestResultEthRequestRPC(*requestRPC, want, got, t, false)
}
