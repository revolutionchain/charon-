package transformer

import (
	"encoding/json"
	"testing"

	"github.com/qtumproject/janus/pkg/eth"
	"github.com/qtumproject/janus/pkg/internal"
	"github.com/qtumproject/janus/pkg/qtum"
)

func TestEstimateGasRequest(t *testing.T) {
	request := eth.CallRequest{
		From: "0x1e6f89d7399081b4f8f8aa1ae2805a5efff2f960",
		To:   "0x1e6f89d7399081b4f8f8aa1ae2805a5efff2f960",
		Data: "0x0",
	}
	requestRaw, err := json.Marshal(&request)
	if err != nil {
		t.Fatal(err)
	}
	requestParamsArray := []json.RawMessage{requestRaw}
	requestRPC, err := internal.PrepareEthRPCRequest(1, requestParamsArray)

	if err != nil {
		t.Fatal(err)
	}

	mockedClientDoer := internal.NewDoerMappedMock()
	qtumClient, err := internal.CreateMockedClient(mockedClientDoer)
	if err != nil {
		t.Fatal(err)
	}

	//preparing responses
	fromHexAddressResponse := qtum.FromHexAddressResponse("0x1e6f89d7399081b4f8f8aa1ae2805a5efff2f960")
	err = mockedClientDoer.AddResponseWithRequestID(2, qtum.MethodFromHexAddress, fromHexAddressResponse)
	if err != nil {
		t.Fatal(err)
	}

	callContractResponse := qtum.CallContractResponse{
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
			GasUsed:  21678,
			Excepted: "None",
		},
	}
	err = mockedClientDoer.AddResponseWithRequestID(1, qtum.MethodCallContract, callContractResponse)
	if err != nil {
		t.Fatal(err)
	}

	//preparing proxy & executing request
	proxyEth := ProxyETHCall{qtumClient}
	proxyEthEstimateGas := ProxyETHEstimateGas{&proxyEth}
	got, jsonErr := proxyEthEstimateGas.Request(requestRPC, nil)
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	want := eth.EstimateGasResponse("0x659d")

	internal.CheckTestResultEthRequestCall(request, &want, got, t, false)
}

func TestEstimateGasRequestExecutionReverted(t *testing.T) {
	request := eth.CallRequest{
		From: "0x1e6f89d7399081b4f8f8aa1ae2805a5efff2f960",
		To:   "0x1e6f89d7399081b4f8f8aa1ae2805a5efff2f960",
		Data: "0x0",
	}
	requestRaw, err := json.Marshal(&request)
	if err != nil {
		t.Fatal(err)
	}
	requestParamsArray := []json.RawMessage{requestRaw}
	requestRPC, err := internal.PrepareEthRPCRequest(1, requestParamsArray)

	if err != nil {
		t.Fatal(err)
	}

	mockedClientDoer := internal.NewDoerMappedMock()
	qtumClient, err := internal.CreateMockedClient(mockedClientDoer)
	if err != nil {
		t.Fatal(err)
	}

	//preparing responses
	fromHexAddressResponse := qtum.FromHexAddressResponse("0x1e6f89d7399081b4f8f8aa1ae2805a5efff2f960")
	err = mockedClientDoer.AddResponseWithRequestID(2, qtum.MethodFromHexAddress, fromHexAddressResponse)
	if err != nil {
		t.Fatal(err)
	}

	callContractResponse := qtum.CallContractResponse{
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
			GasUsed:  21678,
			Excepted: "OutOfGas",
		},
	}
	err = mockedClientDoer.AddResponseWithRequestID(1, qtum.MethodCallContract, callContractResponse)
	if err != nil {
		t.Fatal(err)
	}

	//preparing proxy & executing request
	proxyEth := ProxyETHCall{qtumClient}
	proxyEthEstimateGas := ProxyETHEstimateGas{&proxyEth}

	_, got := proxyEthEstimateGas.Request(requestRPC, nil)

	want := eth.NewCallbackError(ErrExecutionReverted.Error())

	internal.CheckTestResultDefault(want, got, t, false)
}

func TestEstimateGasNonVMRequest(t *testing.T) {
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

	if err != nil {
		t.Fatal(err)
	}

	mockedClientDoer := internal.NewDoerMappedMock()
	qtumClient, err := internal.CreateMockedClient(mockedClientDoer)
	if err != nil {
		t.Fatal(err)
	}

	//preparing responses
	fromHexAddressResponse := qtum.FromHexAddressResponse("0x1e6f89d7399081b4f8f8aa1ae2805a5efff2f960")
	err = mockedClientDoer.AddResponseWithRequestID(2, qtum.MethodFromHexAddress, fromHexAddressResponse)
	if err != nil {
		t.Fatal(err)
	}

	callContractResponse := qtum.CallContractResponse{
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
			GasUsed:  21678,
			Excepted: "None",
		},
	}
	err = mockedClientDoer.AddResponseWithRequestID(1, qtum.MethodCallContract, callContractResponse)
	if err != nil {
		t.Fatal(err)
	}

	//preparing proxy & executing request
	proxyEth := ProxyETHCall{qtumClient}
	proxyEthEstimateGas := ProxyETHEstimateGas{&proxyEth}
	got, jsonErr := proxyEthEstimateGas.Request(requestRPC, nil)
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	want := eth.EstimateGasResponse(NonContractVMGasLimit)

	internal.CheckTestResultEthRequestCall(request, &want, got, t, false)
}
