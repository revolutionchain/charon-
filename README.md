[![Github Build Status](https://github.com/revolutionchain/charon/workflows/Openzeppelin/badge.svg)](https://github.com/revolutionchain/charon/actions)
[![Github Build Status](https://github.com/revolutionchain/charon/workflows/Unit%20tests/badge.svg)](https://github.com/revolutionchain/charon/actions)

# Qtum adapter to Ethereum JSON RPC
Charon is a web3 proxy adapter that can be used as a web3 provider to interact with Qtum. It supports HTTP(s) and websockets and the current version enables self hosting of keys.

# Table of Contents

- [Quick start](#quick-start)
- [Public instances](#public-instances)
- [Requirements](#requirements)
- [Installation](#installation)
  - [SSL](#ssl)
  - [Self-signed SSL](#self-signed-ssl)
- [How to use Charon as a Web3 provider](#how-to-use-charon-as-a-web3-provider)
- [How to add Charon to Metamask](#how-to-add-charon-to-metamask)
- [Truffle support](#truffle-support)
- [Ethers support](#ethers-support)
- [Supported ETH methods](#supported-eth-methods)
- [Websocket ETH methods](#websocket-eth-methods-endpoint-at-)
- [Charon methods](#charon-methods)
- [Development methods](#development-methods)
- [Health checks](#health-checks)
- [Deploying and Interacting with a contract using RPC calls](#deploying-and-interacting-with-a-contract-using-rpc-calls)
  - [Assumption parameters](#assumption-parameters)
  - [Deploy the contract](#deploy-the-contract)
  - [Get the transaction using the hash from previous the result](#get-the-transaction-using-the-hash-from-previous-the-result)
  - [Get the transaction receipt](#get-the-transaction-receipt)
  - [Calling the set method](#calling-the-set-method)
  - [Calling the get method](#calling-the-get-method)
- [Differences between EVM chains](DIFFERENCES.md)

## Quick start
### Public instances
#### You can use public instances if you don't need to use eth_sendTransaction or eth_accounts
Mainnet: https://charon.qiswap.com/api/

Testnet: https://testnet-charon.qiswap.com/api/

Regtest: run it locally with ```make quick-start-regtest```

If you need to use eth_sendTransaction, you are going to have to run your own instance pointing to your own QTUM instance

See [(Beta) QTUM ethers-js library](https://github.com/earlgreytech/qtum-ethers) to generate transactions in the browser so you can use public instances

See [Differences between EVM chains](#differences-between-evm-chains) below

## Requirements

- Golang
- Docker
- linux commands: `make`, `curl`

## Installation

```
$ sudo apt install make git golang docker-compose
# Configure GOPATH if not configured
$ export GOPATH=`go env GOPATH`
$ mkdir -p $GOPATH/src/github.com/revolutionchain && \
  cd $GOPATH/src/github.com/revolutionchain && \
  git clone https://github.com/revolutionchain/charon
$ cd $GOPATH/src/github.com/revolutionchain/charon
# Generate self-signed SSL cert (optional)
# If you do this step, Charon will respond in SSL
# otherwise, Charon will respond unencrypted
$ make docker-configure-https
# Pick a network to quick-start with
$ make quick-start-regtest
$ make quick-start-testnet
$ make quick-start-mainnet
```
This will build the docker image for the local version of Charon as well as spin up two containers:

-   One named `charon` running on port 23889
    
-   Another one named `qtum` running on port 3889
    

`make quick-start` will also fund the tests accounts with QTUM in order for you to start testing and developing locally. Additionally, if you need or want to make changes and or additions to Charon, but don't want to go through the hassle of rebuilding the container, you can run the following command at the project root level:
```
$ make run-charon
# For https
$ make docker-configure-https && make run-charon-https
```
Which will run the most current local version of Charon on port 23888, but without rebuilding the image or the local docker container.

Note that Charon will use the hex address for the test base58 Qtum addresses that belong the the local qtum node, for example:
  - qUbxboqjBRp96j3La8D1RYkyqx5uQbJPoW (hex 0x7926223070547d2d15b2ef5e7383e541c338ffe9 )
  - qLn9vqbr2Gx3TsVR9QyTVB5mrMoh4x43Uf (hex 0x2352be3db3177f0a07efbe6da5857615b8c9901d )

### SSL
SSL keys and certificates go inside the https folder (mounted at `/https` in the container) and use `--https-key` and `--https-cert` parameters. If the specified files do not exist, it will fall back to http.

### Self-signed SSL
To generate self-signed certificates with docker for local development the following script will generate SSL certificates and drop them into the https folder

```
$ make docker-configure-https
```

## How to use Charon as a Web3 provider

Once Charon is successfully running, all one has to do is point your desired framework to Charon in order to use it as your web3 provider. Lets say you want to use truffle for example, in this case all you have to do is go to your truffle-config.js file and add charon as a network:
```
module.exports = {
  networks: {
    charon: {
      host: "127.0.0.1",
      port: 23889,
      network_id: "*",
      gasPrice: "0x5d21dba000"
    },
    ...
  },
...
}
```

## How to add Charon to Metamask

Getting Charon to work with Metamask requires two things
- [Configuring Metamask to point to Charon](metamask)
- Locally signing transactions with a Metamask fork
  - [(Alpha) QTUM Metamask fork](https://github.com/earlgreytech/metamask-extension/releases)

## Truffle support

Hosting your own Charon and blockchain instance works similarly to geth and is supported

Client side transaction signing is supported with [hdwallet-provider](https://www.npmjs.com/package/@qtumproject/hdwallet-provider) underneath it uses [qtum-ethers-wrapper](https://github.com/revolutionchain/qtum-ethers) to construct raw transactions

See [truffle unbox qtumproject/react-box](https://github.com/revolutionchain/react-box) for an example truffle-config file

## Ethers support

Ethers is supported, use [qtum-ethers-wrapper](https://github.com/revolutionchain/qtum-ethers)

## Supported ETH methods

-   [web3_clientVersion](pkg/transformer/web3_clientVersion.go)
-   [web3_sha3](pkg/transformer/web3_sha3.go)
-   [net_version](pkg/transformer/eth_net_version.go)
-   [net_listening](pkg/transformer/eth_net_listening.go)
-   [net_peerCount](pkg/transformer/eth_net_peerCount.go)
-   [eth_protocolVersion](pkg/transformer/eth_protocolVersion.go)
-   [eth_chainId](pkg/transformer/eth_chainId.go)
-   [eth_mining](pkg/transformer/eth_mining.go)
-   [eth_hashrate](pkg/transformer/eth_hashrate.go)
-   [eth_gasPrice](pkg/transformer/eth_gasPrice.go)
-   [eth_accounts](pkg/transformer/eth_accounts.go)
-   [eth_blockNumber](pkg/transformer/eth_blockNumber.go)
-   [eth_getBalance](pkg/transformer/eth_getBalance.go)
-   [eth_getStorageAt](pkg/transformer/eth_getStorageAt.go)
-   [eth_getTransactionCount](pkg/transformer/eth_getTransactionCount.go)
-   [eth_getCode](pkg/transformer/eth_getCode.go)
-   [eth_sign](pkg/transformer/eth_sign.go)
-   [eth_signTransaction](pkg/transformer/eth_signTransaction.go)
-   [eth_sendTransaction](pkg/transformer/eth_sendTransaction.go)
-   [eth_sendRawTransaction](pkg/transformer/eth_sendRawTransaction.go)
-   [eth_call](pkg/transformer/eth_call.go)
-   [eth_estimateGas](pkg/transformer/eth_estimateGas.go)
-   [eth_getBlockByHash](pkg/transformer/eth_getBlockByHash.go)
-   [eth_getBlockByNumber](pkg/transformer/eth_getBlockByNumber.go)
-   [eth_getTransactionByHash](pkg/transformer/eth_getTransactionByHash.go)
-   [eth_getTransactionByBlockHashAndIndex](pkg/transformer/eth_getTransactionByBlockHashAndIndex.go)
-   [eth_getTransactionByBlockNumberAndIndex](pkg/transformer/eth_getTransactionByBlockNumberAndIndex.go)
-   [eth_getTransactionReceipt](pkg/transformer/eth_getTransactionReceipt.go)
-   [eth_getUncleByBlockHashAndIndex](pkg/transformer/eth_getUncleByBlockHashAndIndex.go)
-   [eth_getCompilers](pkg/transformer/eth_getCompilers.go)
-   [eth_newFilter](pkg/transformer/eth_newFilter.go)
-   [eth_newBlockFilter](pkg/transformer/eth_newBlockFilter.go)
-   [eth_uninstallFilter](pkg/transformer/eth_uninstallFilter.go)
-   [eth_getFilterChanges](pkg/transformer/eth_getFilterChanges.go)
-   [eth_getFilterLogs](pkg/transformer/eth_getFilterLogs.go)
-   [eth_getLogs](pkg/transformer/eth_getLogs.go)

## Websocket ETH methods (endpoint at /)

-   (All the above methods)
-   [eth_subscribe](pkg/transformer/eth_subscribe.go) (only 'logs' for now)
-   [eth_unsubscribe](pkg/transformer/eth_unsubscribe.go)

## Charon methods

-   [qtum_getUTXOs](pkg/transformer/qtum_getUTXOs.go)

## Development methods
Use these to speed up development, but don't rely on them in your dapp

-   [dev_gethexaddress](https://docs.qtum.site/en/Qtum-RPC-API/#gethexaddress) Convert Qtum base58 address to hex
-   [dev_fromhexaddress](https://docs.qtum.site/en/Qtum-RPC-API/#fromhexaddress) Convert from hex to Qtum base58 address for the connected network (strip 0x prefix from address when calling this)
-   [dev_generatetoaddress](https://docs.qtum.site/en/Qtum-RPC-API/#generatetoaddress) Mines blocks in regtest (accepts hex/base58 addresses - keep in mind that to use these coins, you must mine 2000 blocks)

## Health checks

There are two health check endpoints, `GET /live` and `GET /ready` they return 200 or 503 depending on health (if they can connect to revod)

## Deploying and Interacting with a contract using RPC calls


### Assumption parameters

Assume that you have a **contract** like this:

```solidity
pragma solidity ^0.4.18;

contract SimpleStore {
  constructor(uint _value) public {
    value = _value;
  }

  function set(uint newValue) public {
    value = newValue;
  }

  function get() public constant returns (uint) {
    return value;
  }

  uint value;
}
```

so that the **bytecode** is

```
solc --optimize --bin contracts/SimpleStore.sol

======= contracts/SimpleStore.sol:SimpleStore =======
Binary:
608060405234801561001057600080fd5b506040516020806100f2833981016040525160005560bf806100336000396000f30060806040526004361060485763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166360fe47b18114604d5780636d4ce63c146064575b600080fd5b348015605857600080fd5b5060626004356088565b005b348015606f57600080fd5b506076608d565b60408051918252519081900360200190f35b600055565b600054905600a165627a7a7230582049a087087e1fc6da0b68ca259d45a2e369efcbb50e93f9b7fa3e198de6402b810029
```

**constructor parameters** is `0000000000000000000000000000000000000000000000000000000000000001`

### Deploy the contract

```
$ curl --header 'Content-Type: application/json' --data \
     '{"id":"10","jsonrpc":"2.0","method":"eth_sendTransaction","params":[{"from":"0x7926223070547d2d15b2ef5e7383e541c338ffe9","gas":"0x6691b7","gasPrice":"0x5d21dba000","data":"0x608060405234801561001057600080fd5b506040516020806100f2833981016040525160005560bf806100336000396000f30060806040526004361060485763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166360fe47b18114604d5780636d4ce63c146064575b600080fd5b348015605857600080fd5b5060626004356088565b005b348015606f57600080fd5b506076608d565b60408051918252519081900360200190f35b600055565b600054905600a165627a7a7230582049a087087e1fc6da0b68ca259d45a2e369efcbb50e93f9b7fa3e198de6402b8100290000000000000000000000000000000000000000000000000000000000000001"}]}' \
     'http://localhost:23889'

{
  "jsonrpc": "2.0",
  "result": "0xa85cacc6143004139fc68808744ea6125ae984454e0ffa6072ac2f2debb0c2e6",
  "id": "10"
}
```

### Get the transaction using the hash from previous the result

```
$ curl --header 'Content-Type: application/json' --data \
     '{"id":"10","jsonrpc":"2.0","method":"eth_getTransactionByHash","params":["0xa85cacc6143004139fc68808744ea6125ae984454e0ffa6072ac2f2debb0c2e6"]}' \
     'localhost:23889'

{
  "jsonrpc":"2.0",
  "result": {
    "blockHash":"0x1e64595e724ea5161c0597d327072074940f519a6fb285ae60e73a4c996b47a4",
    "blockNumber":"0xc9b5",
    "transactionIndex":"0x5",
    "hash":"0xa85cacc6143004139fc68808744ea6125ae984454e0ffa6072ac2f2debb0c2e6",
    "nonce":"0x0",
    "value":"0x0",
    "input":"0x00",
    "from":"0x7926223070547d2d15b2ef5e7383e541c338ffe9",
    "to":"",
    "gas":"0x363639316237",
    "gasPrice":"0x5d21dba000"
  },
  "id":"10"
}
```

### Get the transaction receipt

```
$ curl --header 'Content-Type: application/json' --data \
     '{"id":"10","jsonrpc":"2.0","method":"eth_getTransactionReceipt","params":["0x6da39dc909debf70a536bbc108e2218fd7bce23305ddc00284075df5dfccc21b"]}' \
     'localhost:23889'

{
  "jsonrpc": "2.0",
  "result": {
    "transactionHash": "0xa85cacc6143004139fc68808744ea6125ae984454e0ffa6072ac2f2debb0c2e6",
    "transactionIndex": "0x5",
    "blockHash": "0x1e64595e724ea5161c0597d327072074940f519a6fb285ae60e73a4c996b47a4",
    "from":"0x7926223070547d2d15b2ef5e7383e541c338ffe9"
    "blockNumber": "0xc9b5",
    "cumulativeGasUsed": "0x8c235",
    "gasUsed": "0x1c071",
    "contractAddress": "0x1286595f8683ae074bc026cf0e587177b36842e2",
    "logs": [],
    "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
    "status": "0x1"
  },
  "id": "10"
}
```

### Calling the set method

the ABI code of set method with param '["2"]' is `60fe47b10000000000000000000000000000000000000000000000000000000000000002`

```
$ curl --header 'Content-Type: application/json' --data \
     '{"id":"10","jsonrpc":"2.0","method":"eth_sendTransaction","params":[{"from":"0x7926223070547d2d15b2ef5e7383e541c338ffe9","gas":"0x6691b7","gasPrice":"0x5d21dba000","to":"0x1286595f8683ae074bc026cf0e587177b36842e2","data":"60fe47b10000000000000000000000000000000000000000000000000000000000000002"}]}' \
     'localhost:23889'

{
  "jsonrpc": "2.0",
  "result": "0x51a286c3bc68335274b9fd255e3988918a999608e305475105385f7ccf838339",
  "id": "10"
}
```

### Calling the get method

get method's ABI code is `6d4ce63c`

```
$ curl --header 'Content-Type: application/json' --data \
     '{"id":"10","jsonrpc":"2.0","method":"eth_call","params":[{"from":"0x7926223070547d2d15b2ef5e7383e541c338ffe9","gas":"0x6691b7","gasPrice":"0x5d21dba000","to":"0x1286595f8683ae074bc026cf0e587177b36842e2","data":"6d4ce63c"},"latest"]}' \
     'localhost:23889'

{
  "jsonrpc": "2.0",
  "result": "0x0000000000000000000000000000000000000000000000000000000000000002",
  "id": "10"
}
```

## Future work
- Transparently translate eth_sendRawTransaction from an EVM transaction to a QTUM transaction if the same key is hosted
- Transparently serve blocks by their Ethereum block hash
- Send all QTUM support via eth_sendTransaction
- For eth_subscribe only the 'logs' type is supported at the moment
