<template>
  <div class="hello">
    <div v-if="web3Detected">
      <b-button v-if="revoConnected">Connected to REVO</b-button>
      <b-button v-else-if="connected" v-on:click="connectToRevo()">Connect to REVO</b-button>
      <b-button v-else v-on:click="connectToWeb3()">Connect</b-button>
    </div>
    <b-button v-else>No Web3 detected - Install metamask</b-button>
  </div>
</template>

<script>
let REVOMainnet = {
  chainId: '0x51', // 81
  chainName: 'REVO Mainnet',
  rpcUrls: ['https://charon.qiswap.com/api/'],
  blockExplorerUrls: ['https://revo.info/'],
  iconUrls: [
    'https://revo.info/images/metamask_icon.svg',
    'https://revo.info/images/metamask_icon.png',
  ],
  nativeCurrency: {
    decimals: 18,
    symbol: 'REVO',
  },
};
let REVOTestNet = {
  chainId: '0x22B9', // 8889
  chainName: 'REVO Testnet',
  rpcUrls: ['https://testnet-charon.qiswap.com/api/'],
  blockExplorerUrls: ['https://testnet.revo.info/'],
  iconUrls: [
    'https://revo.info/images/metamask_icon.svg',
    'https://revo.info/images/metamask_icon.png',
  ],
  nativeCurrency: {
    decimals: 18,
    symbol: 'REVO',
  },
};
let REVORegTest = {
  chainId: '0x22BA', // 8890
  chainName: 'REVO Regtest',
  rpcUrls: ['https://localhost:23889'],
  // blockExplorerUrls: ['https://testnet.revo.info/'],
  iconUrls: [
    'https://revo.info/images/metamask_icon.svg',
    'https://revo.info/images/metamask_icon.png',
  ],
  nativeCurrency: {
    decimals: 18,
    symbol: 'REVO',
  },
};
let config = {
  "0x22B8": REVOMainnet,
  "0x22B9": REVOTestNet,
  "0x22BA": REVORegTest,
};

export default {
  name: 'Web3Button',
  props: {
    msg: String,
    connected: Boolean,
    revoConnected: Boolean,
  },
  computed: {
    web3Detected: function() {
      return !!this.Web3;
    },
  },
  methods: {
    getChainId: function() {
      return window.revo.chainId;
    },
    isOnRevoChainId: function() {
      let chainId = this.getChainId();
      return chainId == REVOMainnet.chainId || chainId == REVOTestNet.chainId;
    },
    connectToWeb3: function(){
      if (this.connected) {
        return;
      }
      let self = this;
      window.revo.request({ method: 'eth_requestAccounts' })
        .then(() => {
          console.log("Emitting web3Connected event");
          let revoConnected = self.isOnRevoChainId();
          let currentlyRevoConnected = self.revoConnected;
          self.$emit("web3Connected", true);
          if (currentlyRevoConnected != revoConnected) {
            console.log("ChainID matches REVO, not prompting to add network to web3, already connected.");
            self.$emit("revoConnected", true);
          }
        })
        .catch((e) => {
          console.log("Connecting to web3 failed", arguments, e);
        })
    },
    connectToRevo: function() {
      console.log("Connecting to Revo, current chainID is", this.getChainId());

      let self = this;
      let revoConfig = config[this.getChainId()] || REVOTestNet;
      console.log("Adding network to Metamask", revoConfig);
      window.revo.request({
        method: "wallet_addEthereumChain",
        params: [revoConfig],
      })
        .then(() => {
          self.$emit("revoConnected", true);
        })
        .catch(() => {
          console.log("Adding network failed", arguments);
        })
    },
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
