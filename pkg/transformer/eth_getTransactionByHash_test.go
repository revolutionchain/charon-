package transformer

import (
	"encoding/json"
	"testing"

	"github.com/revolutionchain/charon/pkg/internal"
	"github.com/revolutionchain/charon/pkg/qtum"
	"github.com/shopspring/decimal"
)

type testData struct {
	TxHash   string
	VoutHex  string
	To       string
	From     string
	Input    string
	Gas      string
	GasPrice string
}

func TestGetTransactionByHashRequest(t *testing.T) {
	//preparing request
	requestParams := []json.RawMessage{[]byte(`"0x11e97fa5877c5df349934bafc02da6218038a427e8ed081f048626fa6eb523f5"`)}
	request, err := internal.PrepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}
	mockedClientDoer := internal.NewDoerMappedMock()
	qtumClient, err := internal.CreateMockedClient(mockedClientDoer)

	internal.SetupGetBlockByHashResponses(t, mockedClientDoer)

	//preparing proxy & executing request
	proxyEth := ProxyETHGetTransactionByHash{qtumClient}
	got, JsonErr := proxyEth.Request(request, internal.NewEchoContext())
	if JsonErr != nil {
		t.Fatal(JsonErr)
	}

	want := internal.GetTransactionByHashResponseData

	internal.CheckTestResultEthRequestRPC(*request, &want, got, t, false)
}

func TestGetTransactionByHashRequestWithContractVout(t *testing.T) {

	testsArray := []testData{
		{
			// Using data from https://qtum.info/tx/d20c5c31536e60decf175caf2cbfba980c3678c0f4b201c9b9fa1440102e6451
			// ASM: "4 25548 40 8588b2c50000000000000000000000000000000000000000000000000000000000000000 57946bb437560b13275c32a468c6fd1e0c2cdd48 OP_CALL",
			TxHash:   "0xd20c5c31536e60decf175caf2cbfba980c3678c0f4b201c9b9fa1440102e6451",
			VoutHex:  "540390d003012844095ea7b300000000000000000000000025495b3a87d82e9d7a71b341addfc0d7bb3475c7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1454fefdb5b31164f66ddb68becd7bdd864cacd65bc2",
			Input:    "0x095ea7b300000000000000000000000025495b3a87d82e9d7a71b341addfc0d7bb3475c7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			To:       "0x54fefdb5b31164f66ddb68becd7bdd864cacd65b",
			Gas:      "0x3d090",
			GasPrice: "0x5d21dba000",
		},
		{
			// Edge case taken from openzeppelin tests
			TxHash:   "1664dbafc1dd3c5264209f384b53c569f18b9acad1433a45458e29d46cfbea3e",
			VoutHex:  "0100010001000100142411fd6feb7c148f58101d0cf6e8c8c45af8f219c2",
			Input:    "0x00",
			To:       "0x2411fd6feb7c148f58101d0cf6e8c8c45af8f219",
			Gas:      "0x0",
			GasPrice: "0x0",
		},
	}
	for _, test := range testsArray {
		mockedClientDoer := internal.NewDoerMappedMock()
		qtumClient, _ := internal.CreateMockedClient(mockedClientDoer)
		requestParams := []json.RawMessage{[]byte(`"` + test.TxHash + `"`)}
		request, err := internal.PrepareEthRPCRequest(1, requestParams)
		if err != nil {
			t.Fatal(err)
		}
		internal.SetupGetBlockByHashResponsesWithVouts(
			t,
			[]*qtum.DecodedRawTransactionOutV{
				{
					Value: decimal.Zero,
					N:     0,
					ScriptPubKey: qtum.DecodedRawTransactionScriptPubKey{
						Hex:       test.VoutHex,
						Addresses: []string{},
					},
				},
			},
			mockedClientDoer,
		)
		proxyEth := ProxyETHGetTransactionByHash{qtumClient}
		got, JsonErr := proxyEth.Request(request, internal.NewEchoContext())
		if JsonErr != nil {
			t.Fatal(JsonErr)
		}

		want := internal.GetTransactionByHashResponseData
		want.Input = test.Input
		want.To = test.To
		want.Gas = test.Gas
		want.GasPrice = test.GasPrice

		internal.CheckTestResultEthRequestRPC(*request, &want, got, t, false)
	}
}

// TODO: This test was copied from the above, with the only change being the ASM in the Vout script. However for some reason a bunch of seemingly unrelated field changed in the respose
// For example the gas and gas price field were suddenly non-zero. So something funky is definitely going on here
func TestGetTransactionByHashRequestWithOpSender(t *testing.T) {
	//? Using data from https://qtum.info/tx/0425fa39feed4cd6c93998159901095c147f8b0043823067dc1d25dabf950ac9
	//preparing request
	requestParams := []json.RawMessage{[]byte(`"0x0425fa39feed4cd6c93998159901095c147f8b0043823067dc1d25dabf950ac9"`)}
	request, err := internal.PrepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}
	mockedClientDoer := internal.NewDoerMappedMock()
	qtumClient, err := internal.CreateMockedClient(mockedClientDoer)

	internal.SetupGetBlockByHashResponsesWithVouts(
		t,
		// TODO: Clean this up, refactor
		[]*qtum.DecodedRawTransactionOutV{
			{
				Value: decimal.Zero,
				N:     0,
				ScriptPubKey: qtum.DecodedRawTransactionScriptPubKey{
					// 'ASM' field has no impact in this unit test
					ASM: "1 81e872329e767a0487de7e970992b13b644f1f4f 6b483045022100b83ef90bc808569fb00e29a0f6209d32c1795207c95a554c091401ac8fa8ab920220694b7ec801efd2facea2026d12e8eb5de7689c637f539a620f24c6da8fff235f0121021104b7672c2e08fe321f1bfaffc3768c2777adeedb857b4313ed9d2f15fc8ce4 OP_SENDER 4 55000 40 a9059cbb000000000000000000000000710e94d7f8a5d7a1e5be52bd783370d6e3008a2a0000000000000000000000000000000000000000000000000000000005f5e100 af1ae4e29253ba755c723bca25e883b8deb777b8 OP_CALL",
					Hex: "01011493594441cb5de8b497ad8467d55412c2a0ef36594c6b6a4730440220396b30b7a2f2af482e585473b7575dd2f989f3f3d7cdee55fa34e93f23d5254d022055326cdcab38c58dc3e65c458bfb656cca8340f59534c00ad98b4d4d3303f459012103379c39b6fb2c705db608f98a8fc064f94c66faf894996ca88595487f9ef04a6ec401040390d0030128043d666e8b140000000000000000000000000000000000000086c2",
					// 'Addresses' field has no impact in this unit test
					Addresses: []string{
						"QXeZZ5MsAF5pPrPy47ZFMmtCpg7RExT4mi",
					},
				},
			},
		},
		mockedClientDoer,
	)

	//preparing proxy & executing request
	proxyEth := ProxyETHGetTransactionByHash{qtumClient}
	got, JsonErr := proxyEth.Request(request, internal.NewEchoContext())
	if JsonErr != nil {
		t.Fatal(JsonErr)
	}

	want := internal.GetTransactionByHashResponseData
	want.Input = "0x3d666e8b"
	want.From = "0x93594441cb5de8b497ad8467d55412c2a0ef3659"
	want.To = "0x0000000000000000000000000000000000000086"
	want.Gas = "0x3d090"
	want.GasPrice = "0x5d21dba000"

	internal.CheckTestResultEthRequestRPC(*request, &want, got, t, false)
}

/*
// TODO: Removing this unit test as the transformer computes the "Amount" value (how much QTUM was transferred out) from the MethodDecodeRawTransaction response
// and the way that the balance is calculated cannot return a precision overflow error
func TestGetTransactionByHashRequest_PrecisionOverflow(t *testing.T) {
	//preparing request
	requestParams := []json.RawMessage{[]byte(`"0x11e97fa5877c5df349934bafc02da6218038a427e8ed081f048626fa6eb523f5"`)}
	request, err := prepareEthRPCRequest(1, requestParams)
	if err != nil {
		t.Fatal(err)
	}
	mockedClientDoer := newDoerMappedMock()
	qtumClient, err := createMockedClient(mockedClientDoer)

	//preparing answer to "getblockhash"
	getTransactionResponse := qtum.GetTransactionResponse{
		Amount:            decimal.NewFromFloat(0.20689141234),
		Fee:               decimal.NewFromFloat(-0.2012),
		Confirmations:     2,
		BlockHash:         "ea26fd59a2145dcecd0e2f81b701019b51ca754b6c782114825798973d8187d6",
		BlockIndex:        2,
		BlockTime:         1533092896,
		ID:                "11e97fa5877c5df349934bafc02da6218038a427e8ed081f048626fa6eb523f5",
		Time:              1533092879,
		ReceivedAt:        1533092879,
		Bip125Replaceable: "no",
		Details: []*qtum.TransactionDetail{{Account: "",
			Category:  "send",
			Amount:    decimal.NewFromInt(0),
			Vout:      0,
			Fee:       decimal.NewFromFloat(-0.2012),
			Abandoned: false}},
		Hex: "020000000159c0514feea50f915854d9ec45bc6458bb14419c78b17e7be3f7fd5f563475b5010000006a473044022072d64a1f4ea2d54b7b05050fc853ab192c91cc5ca17e23007867f92f2ab59d9202202b8c9ab9348c8edbb3b98b1788382c8f37642ec9bd6a4429817ab79927319200012103520b1500a400483f19b93c4cb277a2f29693ea9d6739daaf6ae6e971d29e3140feffffff02000000000000000063010403400d0301644440c10f190000000000000000000000006b22910b1e302cf74803ffd1691c2ecb858d3712000000000000000000000000000000000000000000000000000000000000000a14be528c8378ff082e4ba43cb1baa363dbf3f577bfc260e66272970100001976a9146b22910b1e302cf74803ffd1691c2ecb858d371288acb00f0000",
	}
	err = mockedClientDoer.AddResponseWithRequestID(2, qtum.MethodGetTransaction, getTransactionResponse)
	if err != nil {
		t.Fatal(err)
	}

	decodedRawTransactionResponse := qtum.DecodedRawTransactionResponse{
		ID:       "11e97fa5877c5df349934bafc02da6218038a427e8ed081f048626fa6eb523f5",
		Hash:     "d0fe0caa1b798c36da37e9118a06a7d151632d670b82d1c7dc3985577a71880f",
		Size:     552,
		Vsize:    552,
		Version:  2,
		Locktime: 608,
		Vins: []*qtum.DecodedRawTransactionInV{{
			TxID: "7f5350dc474f2953a3f30282c1afcad2fb61cdcea5bd949c808ecc6f64ce1503",
			Vout: 0,
			ScriptSig: struct {
				Asm string `json:"asm"`
				Hex string `json:"hex"`
			}{
				Asm: "3045022100af4de764705dbd3c0c116d73fe0a2b78c3fab6822096ba2907cfdae2bb28784102206304340a6d260b364ef86d6b19f2b75c5e55b89fb2f93ea72c05e09ee037f60b[ALL] 03520b1500a400483f19b93c4cb277a2f29693ea9d6739daaf6ae6e971d29e3140",
				Hex: "483045022100af4de764705dbd3c0c116d73fe0a2b78c3fab6822096ba2907cfdae2bb28784102206304340a6d260b364ef86d6b19f2b75c5e55b89fb2f93ea72c05e09ee037f60b012103520b1500a400483f19b93c4cb277a2f29693ea9d6739daaf6ae6e971d29e3140",
			},
		}},
		Vouts: []*qtum.DecodedRawTransactionOutV{},
	}
	err = mockedClientDoer.AddResponseWithRequestID(3, qtum.MethodDecodeRawTransaction, decodedRawTransactionResponse)
	if err != nil {
		t.Fatal(err)
	}

	getBlockResponse := qtum.GetBlockResponse{
		Hash:              "bba11e1bacc69ba535d478cf1f2e542da3735a517b0b8eebaf7e6bb25eeb48c5",
		Confirmations:     1,
		Strippedsize:      584,
		Size:              620,
		Weight:            2372,
		Height:            3983,
		Version:           536870912,
		VersionHex:        "20000000",
		Merkleroot:        "0b5f03dc9d456c63c587cc554b70c1232449be43d1df62bc25a493b04de90334",
		Time:              1536551888,
		Mediantime:        1536551728,
		Nonce:             0,
		Bits:              "207fffff",
		Difficulty:        4.656542373906925,
		Chainwork:         "0000000000000000000000000000000000000000000000000000000000001f20",
		HashStateRoot:     "3e49216e58f1ad9e6823b5095dc532f0a6cc44943d36ff4a7b1aa474e172d672",
		HashUTXORoot:      "130a3e712d9f8b06b83f5ebf02b27542fb682cdff3ce1af1c17b804729d88a47",
		Previousblockhash: "6d7d56af09383301e1bb32a97d4a5c0661d62302c06a778487d919b7115543be",
		Flags:             "proof-of-stake",
		Proofhash:         "15bd6006ecbab06708f705ecf68664b78b388e4d51416cdafb019d5b90239877",
		Modifier:          "a79c00d1d570743ca8135a173d535258026d26bafbc5f3d951c3d33486b1f120",
		Txs: []string{"3208dc44733cbfa11654ad5651305428de473ef1e61a1ec07b0c1a5f4843be91",
			"8fcd819194cce6a8454b2bec334d3448df4f097e9cdc36707bfd569900268950"},
		Nextblockhash: "d7758774cfdd6bab7774aa891ae035f1dc5a2ff44240784b5e7bdfd43a7a6ec1",
		Signature:     "3045022100a6ab6c2b14b1f73e734f1a61d4d22385748e48836492723a6ab37cdf38525aba022014a51ecb9e51f5a7a851641683541fec6f8f20205d0db49e50b2a4e5daed69d2",
	}
	err = mockedClientDoer.AddResponseWithRequestID(4, qtum.MethodGetBlock, getBlockResponse)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: Get an actual response for this (only addresses are used in this test though)
	getRawTransactionResponse := qtum.GetRawTransactionResponse{
		Vouts: []qtum.RawTransactionVout{
			{
				Details: struct {
					Addresses []string `json:"addresses"`
					Asm       string   `json:"asm"`
					Hex       string   `json:"hex"`
					// ReqSigs   interface{} `json:"reqSigs"`
					Type string `json:"type"`
				}{
					Addresses: []string{
						"7926223070547d2d15b2ef5e7383e541c338ffe9",
					},
				},
			},
		},
	}
	err = mockedClientDoer.AddResponseWithRequestID(4, qtum.MethodGetRawTransaction, &getRawTransactionResponse)
	if err != nil {
		t.Fatal(err)
	}

	//preparing proxy & executing request
	proxyEth := ProxyETHGetTransactionByHash{qtumClient}
	_, err = proxyEth.Request(request, internal.NewEchoContext())

	want := string("decimal.BigInt() was not a success")
	if err.Error() != want {
		t.Errorf(
			"error\ninput: %s\nwanted error: %s\ngot: %s",
			request,
			string(mustMarshalIndent(want, "", "  ")),
			string(mustMarshalIndent(err, "", "  ")),
		)
	}
}
*/
