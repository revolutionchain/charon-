import "core-js/stable"
import "regenerator-runtime/runtime"
import {providers, Contract, ethers} from "ethers"
import {RevoProvider, RevoWallet} from "revo-ethers-wrapper"
import {utils} from "web3"
var $ = require( "jquery" );
import AdoptionArtifact from './Adoption.json'
import Pets from './pets.json'
window.$ = $;
window.jQuery = $;

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
  // rpcUrls: ['https://localhost:23889'],
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
  "0x51": REVOMainnet,
  81: REVOMainnet,
  "0x22B9": REVOTestNet,
  8889: REVOTestNet,
  "0x22BA": REVORegTest,
  8890: REVORegTest,
};
config[REVOMainnet.chainId] = REVOMainnet;
config[REVOTestNet.chainId] = REVOTestNet;
config[REVORegTest.chainId] = REVORegTest;

const metamask = true;
window.App = {
  web3Provider: null,
  contracts: {},
  account: "",

  init: function() {
    // Load pets.
    var petsRow = $('#petsRow');
    var petTemplate = $('#petTemplate');

    for (let i = 0; i < Pets.length; i ++) {
      petTemplate.find('.panel-title').text(Pets[i].name);
      petTemplate.find('img').attr('src', Pets[i].picture);
      petTemplate.find('.pet-breed').text(Pets[i].breed);
      petTemplate.find('.pet-age').text(Pets[i].age);
      petTemplate.find('.pet-location').text(Pets[i].location);
      petTemplate.find('.btn-adopt').attr('pets-id', Pets[i].id);

      petsRow.append(petTemplate.html());
    }

    App.login()
    if (!metamask) {
      return App.initEthers();
    }
    return App.initWeb3();
  },

  getChainId: function() {
    return (window.revo || {}).chainId || 8890;
  },
  isOnRevoChainId: function() {
    let chainId = this.getChainId();
    return chainId == REVOMainnet.chainId ||
        chainId == REVOTestNet.chainId ||
        chainId == REVORegTest.chainId;
  },

  initEthers: function() {
    let revoRpcProvider = new RevoProvider((config[this.getChainId()] || {}).rpcUrls[0]);
    let privKey = "1dd19e1648a23aaf2b3d040454d2569bd7f2cd816cf1b9b430682941a98151df";
    // WIF format
    // let privKey = "cMbgxCJrTYUqgcmiC1berh5DFrtY1KeU4PXZ6NZxgenniF1mXCRk";
    let revoWallet = new RevoWallet(privKey, revoRpcProvider);
    
    window.revoWallet = revoWallet;
    App.account = revoWallet.address
    App.web3Provider = revoWallet;
    return App.initContract();
  },

  initWeb3: function() {
    let self = this;
    let revoConfig = config[this.getChainId()] || REVORegTest;
    console.log("Adding network to Metamask", revoConfig);
    window.revo.request({
      method: "wallet_addEthereumChain",
      params: [revoConfig],
    })
      .then(() => {
        console.log("Successfully connected to revo")
        window.revo.request({ method: 'eth_requestAccounts' })
          .then((accounts) => {
            console.log("Successfully logged into metamask", accounts);
            let revoConnected = self.isOnRevoChainId();
            let currentlyRevoConnected = self.revoConnected;
            if (accounts && accounts.length > 0) {
              App.account = accounts[0];
            }
            if (currentlyRevoConnected != revoConnected) {
              console.log("ChainID matches REVO, not prompting to add network to web3, already connected.");
            }
            let revoRpcProvider = new RevoProvider(REVOTestNet.rpcUrls[0]);
            let revoWallet = new RevoWallet("1dd19e1648a23aaf2b3d040454d2569bd7f2cd816cf1b9b430682941a98151df", revoRpcProvider);
            App.account = revoWallet.address
            if (!metamask) {
              App.web3Provider = revoWallet;
            } else {
              App.web3Provider = new providers.Web3Provider(window.revo);
            }
            
            return App.initContract();
          })
          .catch((e) => {
            console.log("Connecting to web3 failed", e);
          })
      })
      .catch(() => {
        console.log("Adding network failed", arguments);
      })
  },

  initContract: async function() {
    let chainId = utils.hexToNumber(this.getChainId())
    console.log("chainId", chainId)
    const artifacts = AdoptionArtifact.networks[''+chainId];
    if (!artifacts) {
      alert("Contracts are not deployed on chain " + chainId);
      return
    }
    if (!metamask) {
      App.contracts.Adoption = new Contract(artifacts.address, AdoptionArtifact.abi, App.web3Provider)
    } else {
      App.contracts.Adoption = new Contract(artifacts.address, AdoptionArtifact.abi, App.web3Provider.getSigner())
    }
    

    // Set the provider for our contract
    // App.contracts.Adoption.setProvider(App.web3Provider);

    // Use our contract to retrieve and mark the adopted pets
    await App.markAdopted();
    return App.bindEvents();
  },

  bindEvents: function() {
    $(document).on('click', '.btn-adopt', App.handleAdopt);
  },

  markAdopted: function(adopters, account) {
    var adoptionInstance;
    return new Promise((resolve, reject) => {
      let deployed = App.contracts.Adoption.deployed();
      deployed.then(function(instance) {
        adoptionInstance = instance;
        return adoptionInstance.getAdopters.call()
          .then(function(adopters) {
            console.log("Current adopters", adopters)
            for (var i = 0; i < adopters.length; i++) {
              const adopter = adopters[i];
              if (adopter !== '0x0000000000000000000000000000000000000000') {
                $('.panel-pet').eq(i).find('button').text('Adopted').attr('disabled', true);
                $('.panel-pet').eq(i).find('.pet-adopter-container').css('display', 'block');
                let adopterLabel = adopter;
                if (adopter === App.account) {
                  adopterLabel = "You"
                }
                $('.panel-pet').eq(i).find('.pet-adopter-address').text(adopterLabel);
              } else {
                $('.panel-pet').eq(i).find('.pet-adopter-container').css('display', 'none');
              }
            }
            resolve()
            console.log("Successfully marked as adopted")
          }).catch(function(err) {
            console.log(err);
            reject(err)
          });
      }).catch(function(err) {
        console.error(err)
      })
    });
  },

  handleAdopt: function(event) {
    event.preventDefault();

    var petId = parseInt($(event.target).data('id'));

    var adoptionInstance;

    App.contracts.Adoption.deployed().then(function(instance) {
      adoptionInstance = instance;

      return adoptionInstance.adopt(petId/*, {from: App.account}*/);
    }).then(function(result) {
      console.log("Successfully adopted")
      return App.markAdopted();
    }).catch(function(err) {
      console.error("Adoption failed", err)
      console.error(err.message);
    });
  },

  login: function() {
  },

  handleLogout: function() {
    localStorage.removeItem("userWalletAddress");

    App.login();
    App.markAdopted();
  }
};

$(function() {
  $(document).ready(function() {
    App.init();
  });
});
