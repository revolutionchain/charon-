package transformer

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/labstack/echo"
	"github.com/revolutionchain/charon/pkg/eth"
)

type Web3Sha3 struct{}

func (p *Web3Sha3) Method() string {
	return "web3_sha3"
}

func (p *Web3Sha3) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	var err error
	var req eth.Web3Sha3Request
	if err = json.Unmarshal(rawreq.Params, &req); err != nil {
		// TODO: Correct error code?
		return nil, eth.NewInvalidParamsError(err.Error())
	}

	message := req.Message
	var decoded []byte
	// zero length should return "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"
	if len(message) != 0 {
		decoded, err = hexutil.Decode(string(message))
		if err != nil {
			return nil, eth.NewCallbackError("Failed to decode")
		}
	}

	return hexutil.Encode(crypto.Keccak256(decoded)), nil
}
