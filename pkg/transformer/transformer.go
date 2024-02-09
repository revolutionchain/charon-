package transformer

import (
	"github.com/go-kit/kit/log"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/revolutionchain/charon/pkg/eth"
	"github.com/revolutionchain/charon/pkg/notifier"
	"github.com/revolutionchain/charon/pkg/revo"
)

type Transformer struct {
	revoClient   *revo.Revo
	debugMode    bool
	logger       log.Logger
	transformers map[string]ETHProxy
}

// New creates a new Transformer
func New(revoClient *revo.Revo, proxies []ETHProxy, opts ...Option) (*Transformer, error) {
	if revoClient == nil {
		return nil, errors.New("revoClient cannot be nil")
	}

	t := &Transformer{
		revoClient: revoClient,
		logger:     log.NewNopLogger(),
	}

	var err error
	for _, p := range proxies {
		if err = t.Register(p); err != nil {
			return nil, err
		}
	}

	for _, opt := range opts {
		if err := opt(t); err != nil {
			return nil, err
		}
	}

	return t, nil
}

// Register registers an ETHProxy to a Transformer
func (t *Transformer) Register(p ETHProxy) error {
	if t.transformers == nil {
		t.transformers = make(map[string]ETHProxy)
	}

	m := p.Method()
	if _, ok := t.transformers[m]; ok {
		return errors.Errorf("method already exist: %s ", m)
	}

	t.transformers[m] = p

	return nil
}

// Transform takes a Transformer and transforms the request from ETH request and returns the proxy request
func (t *Transformer) Transform(req *eth.JSONRPCRequest, c echo.Context) (interface{}, eth.JSONRPCError) {
	proxy, err := t.getProxy(req.Method)
	if err != nil {
		return nil, err
	}
	resp, err := proxy.Request(req, c)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (t *Transformer) getProxy(method string) (ETHProxy, eth.JSONRPCError) {
	proxy, ok := t.transformers[method]
	if !ok {
		return nil, eth.NewMethodNotFoundError(method)
	}
	return proxy, nil
}

func (t *Transformer) IsDebugEnabled() bool {
	return t.debugMode
}

// DefaultProxies are the default proxy methods made available
func DefaultProxies(revoRPCClient *revo.Revo, agent *notifier.Agent) []ETHProxy {
	filter := eth.NewFilterSimulator()
	getFilterChanges := &ProxyETHGetFilterChanges{Revo: revoRPCClient, filter: filter}
	ethCall := &ProxyETHCall{Revo: revoRPCClient}

	ethProxies := []ETHProxy{
		ethCall,
		&ProxyNetListening{Revo: revoRPCClient},
		&ProxyETHPersonalUnlockAccount{},
		&ProxyETHChainId{Revo: revoRPCClient},
		&ProxyETHBlockNumber{Revo: revoRPCClient},
		&ProxyETHHashrate{Revo: revoRPCClient},
		&ProxyETHMining{Revo: revoRPCClient},
		&ProxyETHNetVersion{Revo: revoRPCClient},
		&ProxyETHGetTransactionByHash{Revo: revoRPCClient},
		&ProxyETHGetTransactionByBlockNumberAndIndex{Revo: revoRPCClient},
		&ProxyETHGetLogs{Revo: revoRPCClient},
		&ProxyETHGetTransactionReceipt{Revo: revoRPCClient},
		&ProxyETHSendTransaction{Revo: revoRPCClient},
		&ProxyETHAccounts{Revo: revoRPCClient},
		&ProxyETHGetCode{Revo: revoRPCClient},

		&ProxyETHNewFilter{Revo: revoRPCClient, filter: filter},
		&ProxyETHNewBlockFilter{Revo: revoRPCClient, filter: filter},
		getFilterChanges,
		&ProxyETHGetFilterLogs{ProxyETHGetFilterChanges: getFilterChanges},
		&ProxyETHUninstallFilter{Revo: revoRPCClient, filter: filter},

		&ProxyETHEstimateGas{ProxyETHCall: ethCall},
		&ProxyETHGetBlockByNumber{Revo: revoRPCClient},
		&ProxyETHGetBlockByHash{Revo: revoRPCClient},
		&ProxyETHGetBalance{Revo: revoRPCClient},
		&ProxyETHGetStorageAt{Revo: revoRPCClient},
		&ETHGetCompilers{},
		&ETHProtocolVersion{},
		&ETHGetUncleByBlockHashAndIndex{},
		&ETHGetUncleCountByBlockHash{},
		&ETHGetUncleCountByBlockNumber{},
		&Web3ClientVersion{},
		&Web3Sha3{},
		&ProxyETHSign{Revo: revoRPCClient},
		&ProxyETHGasPrice{Revo: revoRPCClient},
		&ProxyETHTxCount{Revo: revoRPCClient},
		&ProxyETHSignTransaction{Revo: revoRPCClient},
		&ProxyETHSendRawTransaction{Revo: revoRPCClient},

		&ETHSubscribe{Revo: revoRPCClient, Agent: agent},
		&ETHUnsubscribe{Revo: revoRPCClient, Agent: agent},

		&ProxyREVOGetUTXOs{Revo: revoRPCClient},
		&ProxyREVOGenerateToAddress{Revo: revoRPCClient},

		&ProxyNetPeerCount{Revo: revoRPCClient},
	}

	permittedRevoCalls := []string{
		revo.MethodGetHexAddress,
		revo.MethodFromHexAddress,
	}

	for _, revoMethod := range permittedRevoCalls {
		ethProxies = append(
			ethProxies,
			&ProxyREVOGenericStringArguments{
				Revo:   revoRPCClient,
				prefix: "dev",
				method: revoMethod,
			},
		)
	}

	return ethProxies
}

func SetDebug(debug bool) func(*Transformer) error {
	return func(t *Transformer) error {
		t.debugMode = debug
		return nil
	}
}

func SetLogger(l log.Logger) func(*Transformer) error {
	return func(t *Transformer) error {
		t.logger = log.WithPrefix(l, "component", "transformer")
		return nil
	}
}
