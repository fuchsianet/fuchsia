![alt text](https://fuchsia.network/wp-content/uploads/2019/10/FuchsiaSolidColor.svg "Fuchsia") 

Fuchsia
======
The Official Golang implementation of the Fuchsia high performance, scalable, network protocol featuring our bespoke DxPoS Consensus Algorithm (Delegated Crosschain Proof-of-Stake).

### Want to know more about Fuchsia, check out the following resources:
* [fuchsia.network Website](https://fuchsia.network)

### Our Web wallet and demo page
* [Fuchsia Web Wallet](https://fuchsia.network) (Coming soon...)
* [Fuchsia Demo Page](https://fuchsia.network) (Coming soon...)

### Mainnnet Network Utilities 
* [Mainnet Network Statistics](https://fuchsia.network) (Coming soon...)
* [Testnet Network Statistics](https://fuchsia.network) (Coming soon...)
* [Fuchsia Mainnet Explorer](https://fuchsia.network) (Coming soon...)

### Mainnnet Network Endpoint's
* RPC Gateway to Mainnet: ```(Coming soon...)```
* WSS Gateway to Mainnet: ```(Coming soon...)```


## Prerequisites
For optimal performance the Fuchsia Engineering team recommend running your nodes with at least:
1. Ubuntu 18.04.2 LTS (or any Debian based Dist)
2. 120GB SSD
3. 16 GB Ram or greater
4. Golang >= 1.10.8
5. AWS EC2 (to avoid sync issues/temporary blacklisting)


## Setup
1. You will need to install Golang and add it to your ```$PATH``` variable, We recommend running go version go1.10.8 or greater
2. Clone this repository
    * then ```cd``` into it
    * Run ```make all``` this will create all the fuchsia binaries necessary to run fuchsia within the ```./build/bin``` folder
    * It will be necessary to also add this path to your operating systems ```$PATH``` variable.
3. Next create a folder somewhere on your os and name it appropriately - it will be your data directory holding all the chains data including keystores for whatever accounts you generate. It will be used during genesis creation with the ```--datadir``` flag
4. Next up you will need to get the Genesis file from Fuchsia S3 storage [HERE](https://fuchsia.network) (Coming soon...) or, alternatively, if you wish to develop with our testnet you can also get the testnet genesis from [HERE](https://fuchsia.network) (Coming soon...).
5. Initialise the genesis using the following command: 
    * ```fuchsia --datadir /your/datadir init genesis.json```
    * (Optional) If you wish to create an account at this stage you can also do that using ```fuchsia --datadir /your/datadir account new```  and typing your password when prompted
    * (Optional) Fuchsia's infrastructure team also recommend another step which greatly helps P2P discovery of your node, place a file inside of the --datadir directory within the folder called ```/your/datadir/fuchsia``` and call it called ```static-nodes.json``` use either the Mainnet Bootnodes which can be found [HERE](https://fuchsia.network) (Coming soon..) or the Testnet Bootnodes which can be found [HERE](https://fuchsia.network) (Coming soon..) 
6. The next step we will provide you with a startup script that will start a standard Fuchsia node which can be used as a gateway with which you can connect to locally.
    * (Optional) Another tip from the Fuchsia Engineering team it is important to specify the correct external IP of your machine, to simplify this install the ```dig``` command line utility and use our script below it will automatically populate your external IP for you and vastly help with P2P Discovery.

### Install Dig
```
sudo apt update
sudo apt install dnsutils
```

### Example bash script to start a local node connected to the Fuchsia Mainnet
```
#!/bin/bash
printf "\n============================> Fuchsia Start script for MAINNET local Gateway <====================================\n\n" 
IP="$(dig +short myip.opendns.com @resolver1.opendns.com)" 
fuchsia --datadir=/your/datadir\
    --atmos.ethereum.endpoint "https://mainnet.infura.io"\
    --txpool.releaselimit "3000" \
    --targetgaslimit "21000000" \
    --gasprice "1000000000"\
    --ws\
    --wsaddr="0.0.0.0"\
    --wsapi db,eth,debug,net,txpool,web3\
    --wsorigins "*"\
    --rpc\
    --rpcapi db,eth,debug,net,txpool,web3\
    --rpcaddr "0.0.0.0"\
    --rpcvhosts "*"\
    --rpccorsdomain "*"\
    --cache="2048"\
    --syncmode 'full' \
    --networkid "538"\
    --nat "extip:$IP"\
    console
```

### Example bash script to start a local node connected to the Fuchsia Testnet
```
#!/bin/bash
printf "\n============================> Fuchsia Start script for TESTNET local Gateway <====================================\n\n" 
IP="$(dig +short myip.opendns.com @resolver1.opendns.com)" 
fuchsia --datadir=/your/datadir\
    --atmos.ethereum.endpoint "https://rinkeby.infura.io"\
    --txpool.releaselimit "3000" \
    --targetgaslimit "21000000" \
    --gasprice "1000000000"\
    --ws\
    --wsaddr="0.0.0.0"\
    --wsapi db,eth,debug,net,txpool,web3\
    --wsorigins "*"\
    --rpc\
    --rpcapi db,eth,debug,net,txpool,web3\
    --rpcaddr "0.0.0.0"\
    --rpcvhosts "*"\
    --rpccorsdomain "*"\
    --cache="2048"\
    --syncmode 'full' \
    --networkid "418313827693"\
    --nat "extip:$IP"\
    console
  ```

## License

The fuchsia library (i.e. all code outside of the `cmd` directory) is licensed under the
[GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html), also
included in our repository in the `COPYING.LESSER` file.

The fuchsia binaries (i.e. all code inside of the `cmd` directory) is licensed under the
[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also included
in our repository in the `COPYING` file.
> 1.9.0-stable