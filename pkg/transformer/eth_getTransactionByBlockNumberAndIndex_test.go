package transformer

import (
	"encoding/json"
	"testing"

	"github.com/revolutionchain/charon/pkg/internal"
	"github.com/revolutionchain/charon/pkg/revo"
)

func initializeProxyETHGetTransactionByBlockNumberAndIndex(revoClient *revo.Revo) ETHProxy {
	return &ProxyETHGetTransactionByBlockNumberAndIndex{revoClient}
}

func TestGetTransactionByBlockNumberAndIndex(t *testing.T) {
	testETHProxyRequest(
		t,
		initializeProxyETHGetTransactionByBlockNumberAndIndex,
		[]json.RawMessage{[]byte(`"` + internal.GetTransactionByHashBlockNumberHex + `"`), []byte(`"0x0"`)},
		internal.GetTransactionByHashResponseData,
	)
}
