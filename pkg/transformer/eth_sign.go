package transformer

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/labstack/echo"
	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/revo"
	"github.com/revolutionchain/charon/pkg/utils"
)

// ProxyETHGetLogs implements ETHProxy
type ProxyETHSign struct {
	*revo.Revo
}

func (p *ProxyETHSign) Method() string {
	return "eth_sign"
}

func (p *ProxyETHSign) Request(rawreq *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	var req eth.SignRequest
	if err := unmarshalRequest(rawreq.Params, &req); err != nil {
		p.GetDebugLogger().Log("method", p.Method(), "error", err)
		// TODO: Correct error code?
		return nil, eth.NewInvalidParamsError(err.Error())
	}

	addr := utils.RemoveHexPrefix(req.Account)

	acc := p.Revo.Accounts.FindByHexAddress(addr)
	if acc == nil {
		p.GetDebugLogger().Log("method", p.Method(), "account", addr, "msg", "Unknown account")
		return nil, eth.NewInvalidParamsError(fmt.Sprintf("No such account: %s", addr))
	}

	sig, err := signMessage(acc.PrivKey, req.Message)
	if err != nil {
		p.GetDebugLogger().Log("method", p.Method(), "msg", "Failed to sign message", "error", err)
		return nil, eth.NewCallbackError(err.Error())
	}

	p.GetDebugLogger().Log("method", p.Method(), "msg", "Successfully signed message")

	return eth.SignResponse("0x" + hex.EncodeToString(sig)), nil
}

func signMessage(key *btcec.PrivateKey, msg []byte) ([]byte, error) {
	msghash := chainhash.DoubleHashB(paddedMessage(msg))

	secp256k1 := btcec.S256()

	return btcec.SignCompact(secp256k1, key, msghash, true)
}

var revoSignMessagePrefix = []byte("\u0015Revo Signed Message:\n")

func paddedMessage(msg []byte) []byte {
	var wbuf bytes.Buffer

	wbuf.Write(revoSignMessagePrefix)

	var msglenbuf [binary.MaxVarintLen64]byte
	msglen := binary.PutUvarint(msglenbuf[:], uint64(len(msg)))

	wbuf.Write(msglenbuf[:msglen])
	wbuf.Write(msg)

	return wbuf.Bytes()
}
