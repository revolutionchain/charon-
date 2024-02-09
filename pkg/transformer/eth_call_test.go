package transformer

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/internal"
	"github.com/revolutionchain/charon/pkg/revo"
)

func TestEthCallRequest(t *testing.T) {
	//prepare request
	request := eth.CallRequest{
		From: "0x1e6f89d7399081b4f8f8aa1ae2805a5efff2f960",
		To:   "0x1e6f89d7399081b4f8f8aa1ae2805a5efff2f960",
	}
	requestRaw, err := json.Marshal(&request)
	if err != nil {
		t.Fatal(err)
	}
	requestParamsArray := []json.RawMessage{requestRaw}
	requestRPC, err := internal.PrepareEthRPCRequest(1, requestParamsArray)

	clientDoerMock := internal.NewDoerMappedMock()
	revoClient, err := internal.CreateMockedClient(clientDoerMock)

	//preparing response
	callContractResponse := revo.CallContractResponse{
		Address: "1e6f89d7399081b4f8f8aa1ae2805a5efff2f960",
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
			GasUsed:    21678,
			Excepted:   "None",
			NewAddress: "1e6f89d7399081b4f8f8aa1ae2805a5efff2f960",
			Output:     "0000000000000000000000000000000000000000000000000000000000000001",
		},
		TransactionReceipt: struct {
			StateRoot string        `json:"stateRoot"`
			GasUsed   int           `json:"gasUsed"`
			Bloom     string        `json:"bloom"`
			Log       []interface{} `json:"log"`
		}{
			StateRoot: "d44fc5ad43bae52f01ff7eb4a7bba904ee52aea6c41f337aa29754e57c73fba6",
			GasUsed:   21678,
			Bloom:     "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		},
	}
	err = clientDoerMock.AddResponseWithRequestID(1, revo.MethodCallContract, callContractResponse)
	if err != nil {
		t.Fatal(err)
	}

	fromHexAddressResponse := revo.FromHexAddressResponse("0x1e6f89d7399081b4f8f8aa1ae2805a5efff2f960")
	err = clientDoerMock.AddResponseWithRequestID(2, revo.MethodFromHexAddress, fromHexAddressResponse)
	if err != nil {
		t.Fatal(err)
	}

	//preparing proxy & executing
	proxyEth := ProxyETHCall{revoClient}
	if err != nil {
		t.Fatal(err)
	}

	got, jsonErr := proxyEth.Request(requestRPC, internal.NewEchoContext())
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	want := eth.CallResponse("0x0000000000000000000000000000000000000000000000000000000000000001")

	internal.CheckTestResultEthRequestCall(request, &want, got, t, false)
}

func TestRetry(t *testing.T) {
	//prepare request
	request := eth.CallRequest{
		From: "0x1e6f89d7399081b4f8f8aa1ae2805a5efff2f960",
		To:   "0x1e6f89d7399081b4f8f8aa1ae2805a5efff2f960",
	}
	requestRaw, err := json.Marshal(&request)
	if err != nil {
		t.Fatal(err)
	}
	requestParamsArray := []json.RawMessage{requestRaw}
	requestRPC, err := internal.PrepareEthRPCRequest(1, requestParamsArray)

	clientDoerMock := internal.NewDoerMappedMock()
	revoClient, err := internal.CreateMockedClient(clientDoerMock)

	//preparing response
	callContractResponse := revo.CallContractResponse{
		Address: "1e6f89d7399081b4f8f8aa1ae2805a5efff2f960",
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
			GasUsed:    21678,
			Excepted:   "None",
			NewAddress: "1e6f89d7399081b4f8f8aa1ae2805a5efff2f960",
			Output:     "0000000000000000000000000000000000000000000000000000000000000001",
		},
		TransactionReceipt: struct {
			StateRoot string        `json:"stateRoot"`
			GasUsed   int           `json:"gasUsed"`
			Bloom     string        `json:"bloom"`
			Log       []interface{} `json:"log"`
		}{
			StateRoot: "d44fc5ad43bae52f01ff7eb4a7bba904ee52aea6c41f337aa29754e57c73fba6",
			GasUsed:   21678,
			Bloom:     "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		},
	}

	// return REVO is busy response 4 times
	for i := 0; i < 4; i++ {
		clientDoerMock.AddRawResponse(revo.MethodCallContract, []byte(revo.ErrRevoWorkQueueDepth.Error()))
	}
	// on 5th request, return correct value
	clientDoerMock.AddResponseWithRequestID(1, revo.MethodCallContract, callContractResponse)

	fromHexAddressResponse := revo.FromHexAddressResponse("0x1e6f89d7399081b4f8f8aa1ae2805a5efff2f960")
	err = clientDoerMock.AddResponseWithRequestID(2, revo.MethodFromHexAddress, fromHexAddressResponse)
	if err != nil {
		t.Fatal(err)
	}

	//preparing proxy & executing
	proxyEth := ProxyETHCall{revoClient}
	if err != nil {
		t.Fatal(err)
	}

	before := time.Now()

	got, jsonErr := proxyEth.Request(requestRPC, internal.NewEchoContext())
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	after := time.Now()

	want := eth.CallResponse("0x0000000000000000000000000000000000000000000000000000000000000001")

	internal.CheckTestResultEthRequestCall(request, &want, got, t, false)

	if after.Sub(before) < 2*time.Second {
		t.Errorf("Retrying requests was too quick: %v < 2s", after.Sub(before))
	}
}

func TestEthCallRequestOnUnknownContract(t *testing.T) {
	//prepare request
	request := eth.CallRequest{
		From: "0x1e6f89d7399081b4f8f8aa1ae2805a5efff2f960",
		To:   "0x1e6f89d7399081b4f8f8aa1ae2805a5efff2f960",
	}
	requestRaw, err := json.Marshal(&request)
	if err != nil {
		t.Fatal(err)
	}
	requestParamsArray := []json.RawMessage{requestRaw}
	requestRPC, err := internal.PrepareEthRPCRequest(1, requestParamsArray)

	clientDoerMock := internal.NewDoerMappedMock()
	revoClient, err := internal.CreateMockedClient(clientDoerMock)

	fromHexAddressResponse := revo.FromHexAddressResponse("0x1e6f89d7399081b4f8f8aa1ae2805a5efff2f960")
	err = clientDoerMock.AddResponse(revo.MethodFromHexAddress, fromHexAddressResponse)
	if err != nil {
		t.Fatal(err)
	}

	//preparing error response
	unknownAddressResponse := revo.GetErrorResponse(revo.ErrInvalidAddress)
	err = clientDoerMock.AddError(revo.MethodCallContract, unknownAddressResponse)
	if err != nil {
		t.Fatal(err)
	}

	//preparing proxy & executing
	proxyEth := ProxyETHCall{revoClient}
	if err != nil {
		t.Fatal(err)
	}

	got, jsonErr := proxyEth.Request(requestRPC, internal.NewEchoContext())
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	want := eth.CallResponse("0x")

	internal.CheckTestResultEthRequestCall(request, &want, got, t, false)
}
