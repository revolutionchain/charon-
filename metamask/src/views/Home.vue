<template>
  <div class="home">
    <h4>This is a simple demo on how to add Charon as a chain to Metamask</h4>
    <Web3Button
      @web3Connected="setWeb3Connected"
      @revoConnected="setRevoConnected"
      :connected="web3Connected"
      :revoConnected="revoConnected"
    />
  </div>
</template>

<script>
// @ is an alias to /src
import Web3Button from '@/components/Web3Button.vue'

export default {
  name: 'Home',
  components: {
    Web3Button
  },
  data() {
    return {
      web3Connected: false,
      revoConnected: false,
    }
  },
  methods: {
    setWeb3Connected: function(connected) {
      console.log("setWeb3Connected", connected);
      this.web3Connected = connected;
    },
    setRevoConnected: function(revoConnected) {
      console.log("setRevoConnected", revoConnected);
      console.log("TODO: periodically query chain id to detect if the user changes network")
      this.revoConnected = revoConnected;
    }
  }
}
</script>
