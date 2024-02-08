package eth

import (
	"encoding/json"
	"math/big"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/revolutionchain/charon/pkg/utils"
)

type ETHInt struct {
	*big.Int
}

func (i *ETHInt) Hex() string {
	return hexutil.EncodeBig(i.Int)
}

func (i *ETHInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Int)
}

// UnmarshalJSON needs to be able to parse ETHInt from both hex string or number
func (i *ETHInt) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 {
		return errors.New("data must not be empty")
	}

	isNumber := func(data []byte) bool {
		return data[0] != '"' && data[len(data)-1] != '"'
	}

	if isNumber(data) {
		i.Int, err = bigIntFromNumber(data)
		return err
	}

	i.Int, err = bigIntFromHex(data)
	return err
}

func bigIntFromNumber(data json.RawMessage) (*big.Int, error) {
	var v *big.Int
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, errors.Wrap(err, "json unmarshal")
	}
	return v, nil
}

func bigIntFromHex(data json.RawMessage) (*big.Int, error) {
	var val string

	if err := json.Unmarshal(data, &val); err != nil {
		return nil, errors.Wrap(err, "json unmarshal")
	}

	i, err := utils.DecodeBig(val)
	if err != nil {
		return nil, errors.Wrap(err, "decoding error")
	}
	return i, nil
}
