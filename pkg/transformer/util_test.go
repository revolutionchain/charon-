package transformer

import (
	"fmt"
	"testing"

	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/internal"
	"github.com/revolutionchain/charon/pkg/revo"
	"github.com/revolutionchain/charon/pkg/utils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestEthValueToRevoAmount(t *testing.T) {
	cases := []map[string]interface{}{
		{
			"in":   "0xde0b6b3a7640000",
			"want": decimal.NewFromFloat(1),
		},
		{

			"in":   "0x6f05b59d3b20000",
			"want": decimal.NewFromFloat(0.5),
		},
		{
			"in":   "0x2540be400",
			"want": decimal.NewFromFloat(0.00000001),
		},
		{
			"in":   "0x1",
			"want": decimal.NewFromInt(0),
		},
	}
	for _, c := range cases {
		in := c["in"].(string)
		want := c["want"].(decimal.Decimal)
		got, err := EthValueToRevoAmount(in, MinimumGas)
		if err != nil {
			t.Error(err)
		}

		// TODO: Refactor to use new testing utilities?
		if !got.Equal(want) {
			t.Errorf("in: %s, want: %v, got: %v", in, want, got)
		}
	}
}

func TestRevoValueToEthAmount(t *testing.T) {
	cases := []decimal.Decimal{
		decimal.NewFromFloat(1),
		decimal.NewFromFloat(0.5),
		decimal.NewFromFloat(0.00000001),
		MinimumGas,
	}
	for _, c := range cases {
		in := c
		eth := RevoDecimalValueToETHAmount(in)
		out := EthDecimalValueToRevoAmount(eth)

		// TODO: Refactor to use new testing utilities?
		if !in.Equals(out) {
			t.Errorf("in: %s, eth: %v, revo: %v", in, eth, out)
		}
	}
}

func TestRevoAmountToEthValue(t *testing.T) {
	in, want := decimal.NewFromFloat(0.1), "0x16345785d8a0000"
	got, err := formatRevoAmount(in)
	if err != nil {
		t.Error(err)
	}

	internal.CheckTestResultUnspecifiedInputMarshal(in, want, got, t, false)
}

func TestLowestRevoAmountToEthValue(t *testing.T) {
	in, want := decimal.NewFromFloat(0.00000001), "0x2540be400"
	got, err := formatRevoAmount(in)
	if err != nil {
		t.Error(err)
	}

	internal.CheckTestResultUnspecifiedInputMarshal(in, want, got, t, false)
}

func TestAddressesConversion(t *testing.T) {
	t.Parallel()

	inputs := []struct {
		revoChain   string
		ethAddress  string
		revoAddress string
	}{
		{
			revoChain:   revo.ChainTest,
			ethAddress:  "6c89a1a6ca2ae7c00b248bb2832d6f480f27da68",
			revoAddress: "qTTH1Yr2eKCuDLqfxUyBLCAjmomQ8pyrBt",
		},

		// Test cases for addresses defined here:
		// 	- https://github.com/hayeah/openzeppelin-solidity/blob/revo/REVO-NOTES.md#create-test-accounts
		//
		// NOTE: Ethereum addresses are without `0x` prefix, as it expects by conversion functions
		{
			revoChain:   revo.ChainTest,
			ethAddress:  "7926223070547d2d15b2ef5e7383e541c338ffe9",
			revoAddress: "qUbxboqjBRp96j3La8D1RYkyqx5uQbJPoW",
		},
		{
			revoChain:   revo.ChainTest,
			ethAddress:  "2352be3db3177f0a07efbe6da5857615b8c9901d",
			revoAddress: "qLn9vqbr2Gx3TsVR9QyTVB5mrMoh4x43Uf",
		},
		{
			revoChain:   revo.ChainTest,
			ethAddress:  "69b004ac2b3993bf2fdf56b02746a1f57997420d",
			revoAddress: "qTCCy8qy7pW94EApdoBjYc1vQ2w68UnXPi",
		},
		{
			revoChain:   revo.ChainTest,
			ethAddress:  "8c647515f03daeefd09872d7530fa8d8450f069a",
			revoAddress: "qWMi6ne9mDQFatRGejxdDYVUV9rQVkAFGp",
		},
		{
			revoChain:   revo.ChainTest,
			ethAddress:  "2191744eb5ebeac90e523a817b77a83a0058003b",
			revoAddress: "qLcshhsRS6HKeTKRYFdpXnGVZxw96QQcfm",
		},
		{
			revoChain:   revo.ChainTest,
			ethAddress:  "88b0bf4b301c21f8a47be2188bad6467ad556dcf",
			revoAddress: "qW28njWueNpBXYWj2KDmtFG2gbLeALeHfV",
		},
	}

	for i, in := range inputs {
		var (
			in       = in
			testDesc = fmt.Sprintf("#%d", i)
		)
		// TODO: Investigate why this testing setup is so different
		t.Run(testDesc, func(t *testing.T) {
			revoAddress, err := convertETHAddress(in.ethAddress, in.revoChain)
			require.NoError(t, err, "couldn't convert Ethereum address to Revo address")
			require.Equal(t, in.revoAddress, revoAddress, "unexpected converted Revo address value")

			ethAddress, err := utils.ConvertRevoAddress(in.revoAddress)
			require.NoError(t, err, "couldn't convert Revo address to Ethereum address")
			require.Equal(t, in.ethAddress, ethAddress, "unexpected converted Ethereum address value")
		})
	}
}

func TestSendTransactionRequestHasDefaultGasPriceAndAmount(t *testing.T) {
	var req eth.SendTransactionRequest
	err := unmarshalRequest([]byte(`[{}]`), &req)
	if err != nil {
		t.Fatal(err)
	}
	defaultGasPriceInWei := req.GasPrice.Int
	defaultGasPriceInREVO := EthDecimalValueToRevoAmount(decimal.NewFromBigInt(defaultGasPriceInWei, 1))

	// TODO: Refactor to use new testing utilities?
	if !defaultGasPriceInREVO.Equals(MinimumGas) {
		t.Fatalf("Default gas price does not convert to REVO minimum gas price, got: %s want: %s", defaultGasPriceInREVO.String(), MinimumGas.String())
	}
	if eth.DefaultGasAmountForRevo.String() != req.Gas.Int.String() {
		t.Fatalf("Default gas amount does not match expected default, got: %s want: %s", req.Gas.Int.String(), eth.DefaultGasAmountForRevo.String())
	}
}
