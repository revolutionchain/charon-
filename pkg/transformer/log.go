package transformer

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/revolutionchain/charon/pkg/revo"
)

func GetLogger(proxy ETHProxy, q *revo.Revo) log.Logger {
	method := proxy.Method()
	logger := q.Client.GetLogger()
	return log.WithPrefix(level.Info(logger), method)
}

func GetLoggerFromETHCall(proxy *ProxyETHCall) log.Logger {
	return GetLogger(proxy, proxy.Revo)
}

func GetDebugLogger(proxy ETHProxy, q *revo.Revo) log.Logger {
	method := proxy.Method()
	logger := q.Client.GetDebugLogger()
	return log.WithPrefix(level.Debug(logger), method)
}

func GetDebugLoggerFromETHCall(proxy *ProxyETHCall) log.Logger {
	return GetDebugLogger(proxy, proxy.Revo)
}
