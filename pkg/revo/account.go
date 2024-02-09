package revo

import (
	"encoding/hex"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

type Accounts []*btcutil.WIF

func (as Accounts) FindByHexAddress(addr string) *btcutil.WIF {
	for _, a := range as {
		acc := &Account{a}

		if addr == acc.ToHexAddress() {
			return a
		}
	}

	return nil
}

type Account struct {
	*btcutil.WIF
}

func (a *Account) ToHexAddress() string {
	// wif := (*btcutil.WIF)(a)

	keyid := btcutil.Hash160(a.SerializePubKey())
	return hex.EncodeToString(keyid)
}

var revoMainNetParams = chaincfg.MainNetParams
var revoTestNetParams = chaincfg.MainNetParams

func init() {
	revoMainNetParams.PubKeyHashAddrID = 60
	revoMainNetParams.ScriptHashAddrID = 50

	revoTestNetParams.PubKeyHashAddrID = 65
	revoTestNetParams.ScriptHashAddrID = 50
}

func (a *Account) ToBase58Address(isMain bool) (string, error) {
	params := &revoMainNetParams
	if !isMain {
		params = &revoTestNetParams
	}

	addr, err := btcutil.NewAddressPubKey(a.SerializePubKey(), params)
	if err != nil {
		return "", err
	}

	return addr.AddressPubKeyHash().String(), nil
}
