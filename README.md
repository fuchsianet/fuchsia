![FCH.Network](https://fch.network/wp-content/uploads/2019/10/fhcnetwork.svg "FCH") 

FCH.Network
======
The Official Golang implementation of the FCH high performance, scalable, network protocol featuring our bespoke DxPoS Consensus Algorithm (Delegated Crosschain Proof-of-Stake).

### Want to know more about FCH, check out the following resources:
* [https://fch.network](https://fch.network/) website

### Our Web wallet
* [FCH.Network web wallet](https://wallet.fch.network/)

### Mainnnet Network Utilities 
* [FCH.Network Stats](https://stats.fch.network/)
* [FCH.Network Explorer](https://blockscout.fch.network/)

### Mainnnet Network Endpoint's
* RPC Gateway to Mainnet: ```https://gateway.fch.network```
* WSS Gateway to Mainnet: ```wss://gateway.fch.network/wss```


## Prerequisites
For optimal performance the FCH Engineering team recommend running your nodes with at least:
1. Ubuntu 18.04.2 LTS (or any Debian based Dist)
2. 120GB SSD
3. 16 GB Ram or greater
4. Golang >= 1.10.8
5. AWS EC2 (to avoid sync issues/temporary blacklisting)


## Setup
1. You will need to install Golang and add it to your ```$PATH``` variable, We recommend running go version go1.10.8 or greater
2. Clone this repository
    * then ```cd``` into it
    * Run ```make all``` this will create all the FCH binaries necessary to run FCH within the ```./build/bin``` folder
    * **It will be necessary** to also add this path to your operating systems ```$PATH``` variable.
3. Next create a folder somewhere on your os and name it appropriately - it will be your data directory holding all the chains data including keystores for whatever accounts you generate. It will be used during genesis creation with the ```--datadir``` flag
4. Next up you will need to get the [genesis file](https://github.com/fchnetwork/fch/tree/master/fch/genesis.json).
5. Initialise the genesis using the following command: 
    * ```fch --datadir /your/datadir init genesis.json```
    * (Optional) If you wish to create an account at this stage you can also do that using ```fch --datadir /your/datadir account new```  and typing your password when prompted
    * (Optional) FCH's infrastructure team also recommend another step which greatly helps P2P discovery of your node, place a file inside of the --datadir directory within the folder called ```/your/datadir/fch``` and call it called ```static-nodes.json``` which can be found [here](https://github.com/fchnetwork/fch/tree/master/fch/static-nodes.json)\ 
6. The next step we will provide you with a startup script that will start a standard FCH node which can be used as a gateway with which you can connect to locally.
    * (Optional) Another tip from the FCH Engineering team it is important to specify the correct external IP of your machine, to simplify this install the ```dig``` command line utility and use our script below it will automatically populate your external IP for you and vastly help with P2P Discovery.

### Install Dig
```
sudo apt update
sudo apt install dnsutils
```

### Example bash script to start a local node connected to the FCH.Network Mainnet
```
#!/bin/bash
printf "\n============================> FCH Start script for MAINNET node <====================================\n\n" 
IP="$(dig +short myip.opendns.com @resolver1.opendns.com)" 
fch --datadir=/your/datadir\
    --atmos.ethereum.endpoint "https://ethereum.fch.network"\
    --txpool.releaselimit "3000" \
    --miner.gastarget="126000000"\
    --gasprice "1000000000"\
    --ethstats="Your Name:RYTTC6sG9N6W6xuJ@stats.fch.network/api"\
    --cache="2048"\
    --syncmode 'full' \
    --networkid "4040"\
    --nat "extip:$IP"\
    --nousb\
    console
```

### Delegate setup

To setup delegate you will need to do next:
1. You should have minimum 10M XF for masternode and 100M XF for validator
2. Set allowance for `0xe2b151d2eF8d7D3058E44b6481B06F71d38253c9` address for amount of tokens you would like to bond into delegate
3. Execute next methods on governance contract `0xe2b151d2eF8d7D3058E44b6481B06F71d38253c9` [Governance ABI](https://github.com/fchnetwork/fch/tree/master/fch//abi/governance.json)
    1. `createDelegate` create new delegate. You will have to specify next:
        1. delegate name - can't be changed
        2. signer address - coinbase address of FCH.Network node which will be used for mining
        3. amount of tokens to bond (approved before)
    2. `getAllDelegatesAddresses()` - the last delegate return is your delegate address. You will need it for next step
4. On delegate contract created above (with address obtained) call `launch()` method to start delegate [Delegate ABI](https://github.com/fchnetwork/fch/tree/master/fch//abi/delegate.json).
5. Start node and enable mining there with `miner.start()`. Your delegate will start mining in 20 minutes from launch.

## License

The FCH library (i.e. all code outside of the `cmd` directory) is licensed under the
[GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html), also
included in our repository in the `COPYING.LESSER` file.

The FCH binaries (i.e. all code inside of the `cmd` directory) is licensed under the
[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also included
in our repository in the `COPYING` file.
> 1.9.0-stable
