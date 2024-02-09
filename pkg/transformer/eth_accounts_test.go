package transformer

import (
	"encoding/json"
	"testing"

	"github.com/btcsuite/btcutil"
	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/internal"
	"github.com/revolutionchain/charon/pkg/revo"
)

func TestAccountRequest(t *testing.T) {
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

	exampleAcc1, err := btcutil.DecodeWIF("5JK4Gu9nxCvsCxiq9Zf3KdmA9ACza6dUn5BRLVWAYEtQabdnJ89")
	if err != nil {
		t.Fatal(err)
	}
	exampleAcc2, err := btcutil.DecodeWIF("5JwvXtv6YCa17XNDHJ6CJaveg4mrpqFvcjdrh9FZWZEvGFpUxec")
	if err != nil {
		t.Fatal(err)
	}

	revoClient.Accounts = append(revoClient.Accounts, exampleAcc1, exampleAcc2)

	//preparing proxy & executing request
	proxyEth := ProxyETHAccounts{revoClient}
	got, jsonErr := proxyEth.Request(request, internal.NewEchoContext())
	if jsonErr != nil {
		t.Fatal(jsonErr.Error())
	}

	want := eth.AccountsResponse{"0x6d358cf96533189dd5a602d0937fddf0888ad3ae", "0x7e22630f90e6db16283af2c6b04f688117a55db4"}

	internal.CheckTestResultEthRequestRPC(*request, want, got, t, false)
}

func TestAccountMethod(t *testing.T) {
	mockedClientDoer := internal.NewDoerMappedMock()
	revoClient, err := internal.CreateMockedClient(mockedClientDoer)
	if err != nil {
		t.Fatal(err)
	}
	//preparing proxy & executing request
	proxyEth := ProxyETHAccounts{revoClient}
	got := proxyEth.Method()

	want := string("eth_accounts")

	internal.CheckTestResultDefault(want, got, t, false)
}
func TestAccountToResponse(t *testing.T) {
	mockedClientDoer := internal.NewDoerMappedMock()
	revoClient, err := internal.CreateMockedClient(mockedClientDoer)
	if err != nil {
		t.Fatal(err)
	}
	proxyEth := ProxyETHAccounts{revoClient}
	callResponse := revo.CallContractResponse{
		ExecutionResult: struct {
			GasUsed         int    `json:"gasUsed"`
			Excepted        string `json:"excepted"`
			ExceptedMessage string `json:"exceptedMessage"`
			NewAddress      string `json:"newAddress"`
			Output          string `json:"output"`
			CodeDeposit     int    `json:"codeDeposit"`
			GasRefunded     int    `json:"gasRefunded"`
			DepositSize     int    `json:"depositSize"`
			GasForDeposit   int    `json:"gasForDeposit"`
		}{
			Output: "0x0000000000000000000000000000000000000000000000000000000000000002",
		},
	}

	got := *proxyEth.ToResponse(&callResponse)
	want := eth.CallResponse("0x0000000000000000000000000000000000000000000000000000000000000002")

	internal.CheckTestResultDefault(want, got, t, false)
}
