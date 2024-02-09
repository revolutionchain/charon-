package transformer

import (
	"encoding/json"
	"testing"

	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/internal"
	"github.com/revolutionchain/charon/pkg/revo"
)

func TestGetStorageAtRequestWithNoLeadingZeros(t *testing.T) {
	index := "abcde"
	blockNumber := "0x1234"
	requestParams := []json.RawMessage{[]byte(`"` + internal.GetTransactionByHashBlockNumberHex + `"`), []byte(`"0x` + index + `"`), []byte(`"` + blockNumber + `"`)}
	request, err := internal.PrepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}

	mockedClientDoer := internal.NewDoerMappedMock()
	revoClient, err := internal.CreateMockedClient(mockedClientDoer)

	value := "0x012341231441234123412343211234abcde12342332100000223030004005000"

	getStorageResponse := revo.GetStorageResponse{}
	getStorageResponse[leftPadStringWithZerosTo64Bytes("12345")] = make(map[string]string)
	getStorageResponse[leftPadStringWithZerosTo64Bytes("12345")][leftPadStringWithZerosTo64Bytes(index)] = value
	err = mockedClientDoer.AddResponseWithRequestID(2, revo.MethodGetStorage, getStorageResponse)
	if err != nil {
		t.Fatal(err)
	}

	//preparing proxy & executing request
	proxyEth := ProxyETHGetStorageAt{revoClient}
	got, jsonErr := proxyEth.Request(request, internal.NewEchoContext())
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	want := eth.GetStorageResponse(value)

	internal.CheckTestResultEthRequestRPC(*request, &want, got, t, false)
}

func TestGetStorageAtRequestWithLeadingZeros(t *testing.T) {
	index := leftPadStringWithZerosTo64Bytes("abcde")
	blockNumber := "0x1234"
	requestParams := []json.RawMessage{[]byte(`"` + internal.GetTransactionByHashBlockNumberHex + `"`), []byte(`"0x` + index + `"`), []byte(`"` + blockNumber + `"`)}
	request, err := internal.PrepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}

	mockedClientDoer := internal.NewDoerMappedMock()
	revoClient, err := internal.CreateMockedClient(mockedClientDoer)

	value := "0x012341231441234123412343211234abcde12342332100000223030004005000"

	getStorageResponse := revo.GetStorageResponse{}
	getStorageResponse[leftPadStringWithZerosTo64Bytes("12345")] = make(map[string]string)
	getStorageResponse[leftPadStringWithZerosTo64Bytes("12345")][leftPadStringWithZerosTo64Bytes(index)] = value
	err = mockedClientDoer.AddResponseWithRequestID(2, revo.MethodGetStorage, getStorageResponse)
	if err != nil {
		t.Fatal(err)
	}

	//preparing proxy & executing request
	proxyEth := ProxyETHGetStorageAt{revoClient}
	got, jsonErr := proxyEth.Request(request, internal.NewEchoContext())
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	want := eth.GetStorageResponse(value)

	internal.CheckTestResultEthRequestRPC(*request, &want, got, t, false)
}

func TestGetStorageAtUnknownFieldRequest(t *testing.T) {
	index := "abcde"
	blockNumber := "0x1234"
	requestParams := []json.RawMessage{[]byte(`"` + internal.GetTransactionByHashBlockNumberHex + `"`), []byte(`"0x1234"`), []byte(`"` + blockNumber + `"`)}
	request, err := internal.PrepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}

	mockedClientDoer := internal.NewDoerMappedMock()
	revoClient, err := internal.CreateMockedClient(mockedClientDoer)

	unknownValue := "0x0000000000000000000000000000000000000000000000000000000000000000"
	value := "0x012341231441234123412343211234abcde12342332100000223030004005000"

	getStorageResponse := revo.GetStorageResponse{}
	getStorageResponse[leftPadStringWithZerosTo64Bytes("12345")] = make(map[string]string)
	getStorageResponse[leftPadStringWithZerosTo64Bytes("12345")][leftPadStringWithZerosTo64Bytes(index)] = value
	err = mockedClientDoer.AddResponseWithRequestID(2, revo.MethodGetStorage, getStorageResponse)
	if err != nil {
		t.Fatal(err)
	}

	//preparing proxy & executing request
	proxyEth := ProxyETHGetStorageAt{revoClient}
	got, jsonErr := proxyEth.Request(request, internal.NewEchoContext())
	if jsonErr != nil {
		t.Fatal(jsonErr)
	}

	want := eth.GetStorageResponse(unknownValue)

	internal.CheckTestResultEthRequestRPC(*request, &want, got, t, false)
}

func TestLeftPadStringWithZerosTo64Bytes(t *testing.T) {
	tests := make(map[string]string)

	tests["1"] = "0000000000000000000000000000000000000000000000000000000000000001"
	tests["01"] = "0000000000000000000000000000000000000000000000000000000000000001"
	tests["001"] = "0000000000000000000000000000000000000000000000000000000000000001"
	tests["1001"] = "0000000000000000000000000000000000000000000000000000000000001001"
	tests["0000000000000000000000000000000000000000000000000000000000001001"] = "0000000000000000000000000000000000000000000000000000000000001001"
	tests["1111111111111111111111111111111111111111111111111111111111111111"] = "1111111111111111111111111111111111111111111111111111111111111111"
	tests["21111111111111111111111111111111111111111111111111111111111111111"] = "21111111111111111111111111111111111111111111111111111111111111111"

	for input, expected := range tests {
		result := leftPadStringWithZerosTo64Bytes(input)
		internal.CheckTestResultUnspecifiedInput(input, expected, result, t, false)
	}
}
