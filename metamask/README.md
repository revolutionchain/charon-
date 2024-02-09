# Simple VUE project to switch to REVO network via Metamask

## Project setup
```
npm install
```

### Compiles and hot-reloads for development
```
npm run serve
```

### Compiles and minifies for production
```
npm run build
```

### Customize configuration
See [Configuration Reference](https://cli.vuejs.org/config/).

### wallet_addEthereumChain
```
// request account access
window.revo.request({ method: 'eth_requestAccounts' })
    .then(() => {
        // add chain
        window.revo.request({
            method: "wallet_addEthereumChain",
            params: [{
                {
                    chainId: '0x22B9',
                    chainName: 'Revo Testnet',
                    rpcUrls: ['https://localhost:23889'],
                    blockExplorerUrls: ['https://testnet.revo.info/'],
                    iconUrls: [
                        'https://revo.info/images/metamask_icon.svg',
                        'https://revo.info/images/metamask_icon.png',
                    ],
                    nativeCurrency: {
                        decimals: 18,
                        symbol: 'REVO',
                    },
                }
            }],
        }
    });
```

# Known issues
- Metamask requires https for `rpcUrls` so that must be enabled
  - Either directly through Charon with `--https-key ./path --https-cert ./path2` see [SSL](../README.md#ssl)
  - Through the Makefile `make docker-configure-https && make run-charon-https`
  - Or do it yourself with a proxy (eg, nginx)
