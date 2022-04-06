[![Github Build Status](https://github.com/qtumproject/janus/workflows/Openzeppelin/badge.svg)](https://github.com/qtumproject/janus/actions)
[![Github Build Status](https://github.com/qtumproject/janus/workflows/Unit%20tests/badge.svg)](https://github.com/qtumproject/janus/actions)

# Qtum adapter to Ethereum JSON RPC
Janus is a web3 proxy adapter that can be used as a web3 provider to interact with Qtum. It supports HTTP(s) and websockets and the current version enables self hosting of keys.

# Table of Contents

- [Quick start](#quick-start)
- [Public instances](#public-instances)
- [Requirements](#requirements)
- [Installation](#installation)
  - [SSL](#ssl)
  - [Self-signed SSL](#self-signed-ssl)
- [How to use Janus as a Web3 provider](#how-to-use-janus-as-a-web3-provider)
- [How to add Janus to Metamask](#how-to-add-janus-to-metamask)
- [Supported ETH methods](#supported-eth-methods)
- [Websocket ETH methods](#websocket-eth-methods-endpoint-at-)
- [Janus methods](#janus-methods)
- [Development methods](#development-methods)
- [Health checks](#health-checks)
- [Deploying and Interacting with a contract using RPC calls](#deploying-and-interacting-with-a-contract-using-rpc-calls)
  - [Assumption parameters](#assumption-parameters)
  - [Deploy the contract](#deploy-the-contract)
  - [Get the transaction using the hash from previous the result](#get-the-transaction-using-the-hash-from-previous-the-result)
  - [Get the transaction receipt](#get-the-transaction-receipt)
  - [Calling the set method](#calling-the-set-method)
  - [Calling the get method](#calling-the-get-method)
- [Differences between EVM chains](#differences-between-evm-chains)

## Quick start
### Public instances
#### You can use public instances if you don't need to use eth_sendTransaction or eth_accounts
Mainnet: https://janus.qiswap.com/api/

Testnet: https://testnet-janus.qiswap.com/api/

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
$ mkdir -p $GOPATH/src/github.com/qtumproject && \
  cd $GOPATH/src/github.com/qtumproject && \
  git clone https://github.com/qtumproject/janus
$ cd $GOPATH/src/github.com/qtumproject/janus
# Generate self-signed SSL cert (optional)
# If you do this step, Janus will respond in SSL
# otherwise, Janus will respond unencrypted
$ make docker-configure-https
# Pick a network to quick-start with
$ make quick-start-regtest
$ make quick-start-testnet
$ make quick-start-mainnet
```
This will build the docker image for the local version of Janus as well as spin up two containers:

-   One named `janus` running on port 23889
    
-   Another one named `qtum` running on port 3889
    

`make quick-start` will also fund the tests accounts with QTUM in order for you to start testing and developing locally. Additionally, if you need or want to make changes and or additions to Janus, but don't want to go through the hassle of rebuilding the container, you can run the following command at the project root level:
```
$ make run-janus
# For https
$ make docker-configure-https && make run-janus-https
```
Which will run the most current local version of Janus on port 23888, but without rebuilding the image or the local docker container.

Note that Janus will use the hex address for the test base58 Qtum addresses that belong the the local qtum node, for example:
  - qUbxboqjBRp96j3La8D1RYkyqx5uQbJPoW (hex 0x7926223070547d2d15b2ef5e7383e541c338ffe9 )
  - qLn9vqbr2Gx3TsVR9QyTVB5mrMoh4x43Uf (hex 0x2352be3db3177f0a07efbe6da5857615b8c9901d )

### SSL
SSL keys and certificates go inside the https folder (mounted at `/https` in the container) and use `--https-key` and `--https-cert` parameters. If the specified files do not exist, it will fall back to http.

### Self-signed SSL
To generate self-signed certificates with docker for local development the following script will generate SSL certificates and drop them into the https folder

```
$ make docker-configure-https
```

## How to use Janus as a Web3 provider

Once Janus is successfully running, all one has to do is point your desired framework to Janus in order to use it as your web3 provider. Lets say you want to use truffle for example, in this case all you have to do is go to your truffle-config.js file and add janus as a network:
```
module.exports = {
  networks: {
    janus: {
      host: "127.0.0.1",
      port: 23889,
      network_id: "*",
      gasPrice: "0x64"
    },
    ...
  },
...
}
```

## How to add Janus to Metamask

Getting Janus to work with Metamask requires two things
- [Configuring Metamask to point to Janus](metamask)
- Locally signing transactions with a Metamask fork
  - [(Alpha) QTUM Metamask fork](https://github.com/earlgreytech/metamask-extension/releases)

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

## Janus methods

-   [qtum_getUTXOs](pkg/transformer/qtum_getUTXOs.go)

## Development methods
Use these to speed up development, but don't rely on them in your dapp

-   [dev_gethexaddress](https://docs.qtum.site/en/Qtum-RPC-API/#gethexaddress) Convert Qtum base58 address to hex
-   [dev_fromhexaddress](https://docs.qtum.site/en/Qtum-RPC-API/#fromhexaddress) Convert from hex to Qtum base58 address for the connected network (strip 0x prefix from address when calling this)
-   [dev_generatetoaddress](https://docs.qtum.site/en/Qtum-RPC-API/#generatetoaddress) Mines blocks in regtest (accepts hex/base58 addresses - keep in mind that to use these coins, you must mine 2000 blocks)

## Health checks

There are two health check endpoints, `GET /live` and `GET /ready` they return 200 or 503 depending on health (if they can connect to qtumd)

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
     '{"id":"10","jsonrpc":"2.0","method":"eth_sendTransaction","params":[{"from":"0x7926223070547d2d15b2ef5e7383e541c338ffe9","gas":"0x6691b7","gasPrice":"0x64","data":"0x608060405234801561001057600080fd5b506040516020806100f2833981016040525160005560bf806100336000396000f30060806040526004361060485763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166360fe47b18114604d5780636d4ce63c146064575b600080fd5b348015605857600080fd5b5060626004356088565b005b348015606f57600080fd5b506076608d565b60408051918252519081900360200190f35b600055565b600054905600a165627a7a7230582049a087087e1fc6da0b68ca259d45a2e369efcbb50e93f9b7fa3e198de6402b8100290000000000000000000000000000000000000000000000000000000000000001"}]}' \
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
    "gasPrice":"0x3634"
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
     '{"id":"10","jsonrpc":"2.0","method":"eth_sendTransaction","params":[{"from":"0x7926223070547d2d15b2ef5e7383e541c338ffe9","gas":"0x6691b7","gasPrice":"0x64","to":"0x1286595f8683ae074bc026cf0e587177b36842e2","data":"60fe47b10000000000000000000000000000000000000000000000000000000000000002"}]}' \
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
     '{"id":"10","jsonrpc":"2.0","method":"eth_call","params":[{"from":"0x7926223070547d2d15b2ef5e7383e541c338ffe9","gas":"0x6691b7","gasPrice":"0x64","to":"0x1286595f8683ae074bc026cf0e587177b36842e2","data":"6d4ce63c"},"latest"]}' \
     'localhost:23889'

{
  "jsonrpc": "2.0",
  "result": "0x0000000000000000000000000000000000000000000000000000000000000002",
  "id": "10"
}
```

## Differences between EVM chains
- Transaction signing is incompatible
  - QTUM is based on Bitcoin and therefore requires Bitcoin transaction signing
    - EVM transactions are done with special opcodes in Bitcoin output scripts (OP_CALL/OP_CREATE)
  - Use [(Beta) QTUM ethers-js library](https://github.com/earlgreytech/qtum-ethers) to sign transactions for use in eth_sendRawTransaction
    - Currently, the library only supports sending 1 tx per block due to Bitcoin inputs being re-used so test your code to redo transactions if they are rejected with eth_sendRawTransaction
      - This will be fixed in a future version
- Solidity
  - msg.value is denoted in satoshis, not wei, your dapp needs to handle this correctly
  - ecrecover won't recover a QTUM address from eth_sign, you will need to implement [QIP6 - btc_ecrecover](https://blog.qtum.org/qip-6-87e7a9743e14) in your contracts
    - [(Beta) QTUM ethers-js library](https://github.com/earlgreytech/qtum-ethers) doesn't currently sign messages correctly that btc_ecrecover can recover
- Sending coins with the creation of a contract will cause a loss of coins
  - This is a Qtum intentional deisgn decision and will not change
  - Janus will prevent this with eth_sendTransaction but will permit it with eth_sendRawTransaction
  - [(Beta) QTUM ethers-js library](https://github.com/earlgreytech/qtum-ethers) will reject creating such a transaction
- Contract address generation differs from EVM chains
  - on EVM chains, the contract address is generated via a hash of the deployer address + the nonce
  - QTUM has no concept of a nonce because it is built on Bitcoin
    - instead the contract address is generated via a hash of the transaction which will always be different because the Bitcoin inputs will be different
    - so, if your app depends on a consistent contract address between deployments on different chains you need to pay special attention to this
    - For contract address generation code, see [generateContractAddress](https://github.com/earlgreytech/qtum-ethers/blob/main/src/lib/helpers/utils.ts)
- Account address generation differs from EVM chains
  - You really only need to worry about this if you need to use the same account address on different chains
  - [eth_accounts](pkg/transformer/eth_accounts.go) and [(Beta) QTUM ethers-js library](https://github.com/earlgreytech/qtum-ethers) will abstract this away from you
  - For account address generation code, see [computeAddress](https://github.com/earlgreytech/qtum-ethers/blob/main/src/lib/helpers/utils.ts)
- Block hash is computed differently from EVM chains
  - If you are generating the blockhash from the block header, it will be wrong
    - we plan to add a compatiblity layer in Janus to transparently serve the correct block when requesting an Ethereum block hash
      - this will eventually require hooking up Janus to a database to keep a map of hash(block header) => QTUM block hash
- Remix
  - Debug calls are not supported so you will not be able to do any debugging in Remix
  - You can use Remix with Janus or [(Alpha) QTUM Metamask fork](https://github.com/earlgreytech/metamask-extension/releases)
- It is possible for a QTUM transaction to have more than one EVM transaction in it
  - this is because QTUM does EVM transactions inside Bitcoin outputs which there can be multiple of
  - Janus and [(Beta) QTUM ethers-js library](https://github.com/earlgreytech/qtum-ethers) will not generate such a transaction
  - Janus will try and work around such a transaction when requesting information about the transaction
    - but it is not possible to map to the EVM model perfectly so there will always be some data missing for these transactions
- QTUM is proof of stake and requires coins to be mature (older than 2000 blocks) to be used in a transaction
  - this includes staking rewards, gas refunds and block rewards on your local regtest environment
    - a gas refund is an output generated by the miner for every EVM transaction in the same block as the EVM transaction takes place in for unused gas
  - Janus will try to use mature coins first and will fall back to immature coins if there are no mature coins left
    - this can result in transactions being rejected
  - [(Beta) QTUM ethers-js library](https://github.com/earlgreytech/qtum-ethers) will not use immature coins for transactions, but if you end up using high gas limits for your transactions you could quickly run out of usable coins
    - if there are no mature coins, the transaction will fail locally
- Bitcoin input scripts
  - Bitcoin has many different types of scripts
    - For a detailed primer on this topic see [A breakdown of Bitcoin "standard" script types (crazy long)](https://www.reddit.com/r/Bitcoin/comments/jmiko9/a_breakdown_of_bitcoin_standard_script_types/)
  - [eth_sendTransaction](/pkg/transformer/eth_sendTransaction.go) delegates transaction signing to QTUM so most input scripts should be supported
  - [(Beta) QTUM ethers-js library](https://github.com/earlgreytech/qtum-ethers) deals with signing transactions locally and only supports Pay to public key hash (P2PKH) scripts, other script types will be ignored and not selected.
    - This can result in your spendable balance being lower than your actual balance.
    - Support for Pay to public key (P2PK) input scripts is on the roadmap
- [eth_estimateGas](/pkg/transformer/eth_estimateGas.go)
  - Gas estimation on QTUM is not perfect, so a buffer of 10% is added in Janus
  - Gas will be refunded in the block that your transaction is mined
    - Keep in mind that to re-use this gas refund, you must wait 2000 blocks
- [eth_sendTransaction](/pkg/transformer/eth_sendTransaction.go)
  - When trying to send all your QTUM Balance in a transaction, in EVM you would do value = total - (gas limit * gas price)
  - Since QTUM uses Bitcoin transactions, the cost of a transaction differs based on how many bytes are in the transaction
    - This means if you have many inputs in a transaction, it will cost more to send
  - There is no easy way to send your entire QTUM balance in a single transaction with Janus
    - However, [(Beta) QTUM ethers-js library](https://github.com/earlgreytech/qtum-ethers) supports this via value = total - (gas limit * gas price)
    - Adding this to Janus is on the roadmap
- Since QTUM runs on Bitcoin, QTUM has the concept of [dust](https://en.bitcoinwiki.org/wiki/Cryptocurrency_dust)
  - Janus delegates transaction signing to QTUM so QTUM will handle dealing with dust
  - [(Beta) QTUM ethers-js library](https://github.com/earlgreytech/qtum-ethers) currently uses dust, but at some point will prevent spending dust by default with a semver change
- On a transfer of Qtum to a Qtum address, there is no receipt generated for such a transfer
- When converting from WEI -> QTUM, precision is lost due to QTUM's smallest demonination being 1 satoshi.
  - 1 satoshi = 0.00000001 QTUM = 10000000000 wei
- QTUM's minimum gas price is 40 satoshi
  - When specifying a gas price in wei lower than that, the minimum gas price will be used (40 satoshi)
  - With the minimum fee per byte being 4 satoshi
- QTUM will reject transactions with very large fees (to prevent accidents)

## Future work
- Transparently translate eth_sendRawTransaction from an EVM transaction to a QTUM transaction if the same key is hosted
- Transparently serve blocks by their Ethereum block hash
- Send all QTUM support via eth_sendTransaction
- For eth_subscribe only the 'logs' type is supported at the moment
