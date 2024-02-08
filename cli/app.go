package cli

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/btcsuite/btcutil"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"github.com/revolutionchain/charon/pkg/analytics"
	"github.com/revolutionchain/charon/pkg/notifier"
	"github.com/revolutionchain/charon/pkg/params"
	"github.com/revolutionchain/charon/pkg/qtum"
	"github.com/revolutionchain/charon/pkg/server"
	"github.com/revolutionchain/charon/pkg/transformer"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("charon", "Qtum adapter to Ethereum JSON RPC")

	accountsFile = app.Flag("accounts", "account private keys (in WIF) returned by eth_accounts").Envar("ACCOUNTS").File()

	qtumRPC             = app.Flag("qtum-rpc", "URL of qtum RPC service").Envar("QTUM_RPC").Default("").String()
	qtumNetwork         = app.Flag("qtum-network", "if 'regtest' (or connected to a regtest node with 'auto') Charon will generate blocks").Envar("QTUM_NETWORK").Default("auto").String()
	generateToAddressTo = app.Flag("generateToAddressTo", "[regtest only] configure address to mine blocks to when mining new transactions in blocks").Envar("GENERATE_TO_ADDRESS").Default("").String()
	bind                = app.Flag("bind", "network interface to bind to (e.g. 0.0.0.0) ").Default("localhost").String()
	port                = app.Flag("port", "port to serve proxy").Default("23889").Int()
	httpsKey            = app.Flag("https-key", "https keyfile").Default("").String()
	httpsCert           = app.Flag("https-cert", "https certificate").Default("").String()
	logFile             = app.Flag("log-file", "write logs to a file").Envar("LOG_FILE").Default("").String()
	matureBlockHeight   = app.Flag("mature-block-height-override", "override how old a coinbase/coinstake needs to be to be considered mature enough for spending (QTUM uses 2000 blocks after the 32s block fork) - if this value is incorrect transactions can be rejected").Int()
	healthCheckPercent  = app.Flag("health-check-healthy-request-amount", "configure the minimum request success rate for healthcheck").Envar("HEALTH_CHECK_REQUEST_PERCENT").Default("80").Int()

	sqlHost     = app.Flag("sql-host", "database hostname").Envar("SQL_HOST").Default("127.0.0.1").String()
	sqlPort     = app.Flag("sql-port", "database port").Envar("SQL_PORT").Default("5432").Int()
	sqlUser     = app.Flag("sql-user", "database username").Envar("SQL_USER").Default("postgres").String()
	sqlPassword = app.Flag("sql-password", "database password").Envar("SQL_PASSWORD").Default("dbpass").String()
	sqlSSL      = app.Flag("sql-ssl", "use SSL to connect to database").Envar("SQL_SSL").Bool()
	sqlDbname   = app.Flag("sql-dbname", "database name").Envar("SQL_DBNAME").Default("postgres").String()

	dbConnectionString = app.Flag("dbstring", "database connection string").String()

	devMode        = app.Flag("dev", "[Insecure] Developer mode").Envar("DEV").Default("false").Bool()
	singleThreaded = app.Flag("singleThreaded", "[Non-production] Process RPC requests in a single thread").Envar("SINGLE_THREADED").Default("false").Bool()

	ignoreUnknownTransactions = app.Flag("ignoreTransactions", "[Development] Ignore transactions inside blocks we can't fetch and return responses instead of failing").Default("false").Bool()
	disableSnipping           = app.Flag("disableSnipping", "[Development] Disable ...snip... in logs").Default("false").Bool()
	hideQtumdLogs             = app.Flag("hideQtumdLogs", "[Development] Hide QTUMD debug logs").Envar("HIDE_QTUMD_LOGS").Default("false").Bool()
)

func loadAccounts(r io.Reader, l log.Logger) qtum.Accounts {
	var accounts qtum.Accounts

	if accountsFile != nil {
		s := bufio.NewScanner(*accountsFile)
		for s.Scan() {
			line := s.Text()

			wif, err := btcutil.DecodeWIF(line)
			if err != nil {
				level.Error(l).Log("msg", "Failed to parse account", "err", err.Error())
				continue
			}

			accounts = append(accounts, wif)
		}
	}

	if len(accounts) > 0 {
		level.Info(l).Log("msg", fmt.Sprintf("Loaded %d accounts", len(accounts)))
	} else {
		level.Warn(l).Log("msg", "No accounts loaded from account file")
	}

	return accounts
}

func action(pc *kingpin.ParseContext) error {
	addr := fmt.Sprintf("%s:%d", *bind, *port)
	writers := []io.Writer{os.Stdout}

	if logFile != nil && (*logFile) != "" {
		_, err := os.Stat(*logFile)
		if os.IsNotExist(err) {
			newLogFile, err := os.Create(*logFile)
			if err != nil {
				return errors.Wrapf(err, "Failed to create log file %s", *logFile)
			} else {
				writers = append(writers, newLogFile)
			}
		} else {
			existingLogFile, err := os.Open(*logFile)
			if err != nil {
				return errors.Wrapf(err, "Failed to open log file %s", *logFile)
			} else {
				writers = append(writers, existingLogFile)
			}
		}
	}

	logWriter := io.MultiWriter(writers...)
	logger := log.NewLogfmtLogger(logWriter)

	if !*devMode {
		logger = level.NewFilter(logger, level.AllowWarn())
	}

	var accounts qtum.Accounts
	if *accountsFile != nil {
		accounts = loadAccounts(*accountsFile, logger)
		(*accountsFile).Close()
	}

	isMain := *qtumNetwork == qtum.ChainMain

	ctx, shutdownQtum := context.WithCancel(context.Background())
	defer shutdownQtum()

	qtumRequestAnalytics := analytics.NewAnalytics(50)

	qtumJSONRPC, err := qtum.NewClient(
		isMain,
		*qtumRPC,
		qtum.SetDebug(*devMode),
		qtum.SetLogWriter(logWriter),
		qtum.SetLogger(logger),
		qtum.SetAccounts(accounts),
		qtum.SetGenerateToAddress(*generateToAddressTo),
		qtum.SetIgnoreUnknownTransactions(*ignoreUnknownTransactions),
		qtum.SetDisableSnippingQtumRpcOutput(*disableSnipping),
		qtum.SetHideQtumdLogs(*hideQtumdLogs),
		qtum.SetMatureBlockHeight(matureBlockHeight),
		qtum.SetContext(ctx),
		qtum.SetSqlHost(*sqlHost),
		qtum.SetSqlPort(*sqlPort),
		qtum.SetSqlUser(*sqlUser),
		qtum.SetSqlPassword(*sqlPassword),
		qtum.SetSqlSSL(*sqlSSL),
		qtum.SetSqlDatabaseName(*sqlDbname),
		qtum.SetSqlConnectionString(*dbConnectionString),
		qtum.SetAnalytics(qtumRequestAnalytics),
	)
	if err != nil {
		return errors.Wrap(err, "Failed to setup QTUM client")
	}

	qtumClient, err := qtum.New(qtumJSONRPC, *qtumNetwork)
	if err != nil {
		return errors.Wrap(err, "Failed to setup QTUM chain")
	}

	agent := notifier.NewAgent(context.Background(), qtumClient, nil)
	proxies := transformer.DefaultProxies(qtumClient, agent)
	t, err := transformer.New(
		qtumClient,
		proxies,
		transformer.SetDebug(*devMode),
		transformer.SetLogger(logger),
	)
	if err != nil {
		return errors.Wrap(err, "transformer#New")
	}
	agent.SetTransformer(t)

	httpsKeyFile := getEmptyStringIfFileDoesntExist(*httpsKey, logger)
	httpsCertFile := getEmptyStringIfFileDoesntExist(*httpsCert, logger)

	s, err := server.New(
		qtumClient,
		t,
		addr,
		server.SetLogWriter(logWriter),
		server.SetLogger(logger),
		server.SetDebug(*devMode),
		server.SetSingleThreaded(*singleThreaded),
		server.SetHttps(httpsKeyFile, httpsCertFile),
		server.SetQtumAnalytics(qtumRequestAnalytics),
		server.SetHealthCheckPercent(healthCheckPercent),
	)
	if err != nil {
		return errors.Wrap(err, "server#New")
	}

	return s.Start()
}

func getEmptyStringIfFileDoesntExist(file string, l log.Logger) string {
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		l.Log("file does not exist", file)
		return ""
	}
	return file
}

func Run() {
	app.Version(params.VersionWithGitSha)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}

func init() {
	app.Action(action)
}
