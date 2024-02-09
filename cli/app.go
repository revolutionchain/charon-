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
	"github.com/revolutionchain/charon/pkg/revo"
	"github.com/revolutionchain/charon/pkg/server"
	"github.com/revolutionchain/charon/pkg/transformer"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("charon", "Revo adapter to Ethereum JSON RPC")

	accountsFile = app.Flag("accounts", "account private keys (in WIF) returned by eth_accounts").Envar("ACCOUNTS").File()

	revoRPC             = app.Flag("revo-rpc", "URL of revo RPC service").Envar("REVO_RPC").Default("").String()
	revoNetwork         = app.Flag("revo-network", "if 'regtest' (or connected to a regtest node with 'auto') Charon will generate blocks").Envar("REVO_NETWORK").Default("auto").String()
	generateToAddressTo = app.Flag("generateToAddressTo", "[regtest only] configure address to mine blocks to when mining new transactions in blocks").Envar("GENERATE_TO_ADDRESS").Default("").String()
	bind                = app.Flag("bind", "network interface to bind to (e.g. 0.0.0.0) ").Default("localhost").String()
	port                = app.Flag("port", "port to serve proxy").Default("23889").Int()
	httpsKey            = app.Flag("https-key", "https keyfile").Default("").String()
	httpsCert           = app.Flag("https-cert", "https certificate").Default("").String()
	logFile             = app.Flag("log-file", "write logs to a file").Envar("LOG_FILE").Default("").String()
	matureBlockHeight   = app.Flag("mature-block-height-override", "override how old a coinbase/coinstake needs to be to be considered mature enough for spending (REVO uses 2000 blocks after the 32s block fork) - if this value is incorrect transactions can be rejected").Int()
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
	hideRevodLogs             = app.Flag("hideRevodLogs", "[Development] Hide REVOD debug logs").Envar("HIDE_REVOD_LOGS").Default("false").Bool()
)

func loadAccounts(r io.Reader, l log.Logger) revo.Accounts {
	var accounts revo.Accounts

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

	var accounts revo.Accounts
	if *accountsFile != nil {
		accounts = loadAccounts(*accountsFile, logger)
		(*accountsFile).Close()
	}

	isMain := *revoNetwork == revo.ChainMain

	ctx, shutdownRevo := context.WithCancel(context.Background())
	defer shutdownRevo()

	revoRequestAnalytics := analytics.NewAnalytics(50)

	revoJSONRPC, err := revo.NewClient(
		isMain,
		*revoRPC,
		revo.SetDebug(*devMode),
		revo.SetLogWriter(logWriter),
		revo.SetLogger(logger),
		revo.SetAccounts(accounts),
		revo.SetGenerateToAddress(*generateToAddressTo),
		revo.SetIgnoreUnknownTransactions(*ignoreUnknownTransactions),
		revo.SetDisableSnippingRevoRpcOutput(*disableSnipping),
		revo.SetHideRevodLogs(*hideRevodLogs),
		revo.SetMatureBlockHeight(matureBlockHeight),
		revo.SetContext(ctx),
		revo.SetSqlHost(*sqlHost),
		revo.SetSqlPort(*sqlPort),
		revo.SetSqlUser(*sqlUser),
		revo.SetSqlPassword(*sqlPassword),
		revo.SetSqlSSL(*sqlSSL),
		revo.SetSqlDatabaseName(*sqlDbname),
		revo.SetSqlConnectionString(*dbConnectionString),
		revo.SetAnalytics(revoRequestAnalytics),
	)
	if err != nil {
		return errors.Wrap(err, "Failed to setup REVO client")
	}

	revoClient, err := revo.New(revoJSONRPC, *revoNetwork)
	if err != nil {
		return errors.Wrap(err, "Failed to setup REVO chain")
	}

	agent := notifier.NewAgent(context.Background(), revoClient, nil)
	proxies := transformer.DefaultProxies(revoClient, agent)
	t, err := transformer.New(
		revoClient,
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
		revoClient,
		t,
		addr,
		server.SetLogWriter(logWriter),
		server.SetLogger(logger),
		server.SetDebug(*devMode),
		server.SetSingleThreaded(*singleThreaded),
		server.SetHttps(httpsKeyFile, httpsCertFile),
		server.SetRevoAnalytics(revoRequestAnalytics),
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
