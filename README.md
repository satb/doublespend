This project simulates a double spend on a local 2 node ethereum network. The setup and the simulation and the criteria for determining the doublespend are shown below.

* Setup
  * Download the ethereum geth repo and build to get the geth binary https://github.com/ethereum/go-ethereum/wiki/Installing-Geth#build-it-from-source-code
  * Create a new genesis file like so and save it to /tmp/genesis.json (or whatever other folder):
    * `{
            "config": {
                "chainId": 15,
                "homesteadBlock": 0,
                "eip150Block": 0,
                "eip155Block": 0,
                "eip158Block": 0
            },
            "difficulty": "100",
            "gasLimit": "2100000",
            "alloc": {
                "7df9a875a174b3bc565e6424a0050ebc1b2d1d82": { "balance": "300000" },
                "f41c74c9ae680c1aa78f42e5647a62f353b7bdde": { "balance": "400000" }
            }
        }`
  * Create the data directories to hold the data for the ethereum nodes
    * `mkdir -p /tmp/eth/60/01 && mkdir -p /tmp/eth/60/02`

* Running two ethereum nodes locally
  * Bootstrap the ethereum nodes with the genesis file created above like so:
    * `./build/bin/geth --identity "node1" --nodiscover --datadir /tmp/eth/60/01 init /tmp/genesis.json`
    * `./build/bin/geth --identity "node2" --nodiscover --datadir /tmp/eth/60/02 init /tmp/genesis.json`
  * Now you are ready to run the local ethereum nodes
    *  `./build/bin/geth --datadir="/tmp/eth/60/01" --networkid 15 --nodiscover --ws  --wsaddr 127.0.0.1 --wsport 4101 --wsapi "eth,net,web3,admin,shh" --rpc --rpcapi 'personal,db,eth,net,web3,admin,miner,txpool' --ipcdisable --port 30301 --rpcaddr 127.0.0.1 --rpcport 8101  --allow-insecure-unlock console 2>> /tmp/eth/60/01.log`
    * `./build/bin/geth --datadir="/tmp/eth/60/02" --networkid 15 --nodiscover --ws  --wsaddr 127.0.0.1 --wsport 4102 --wsapi "eth,net,web3,admin,shh" --wsorigins "*"  --rpc --rpcapi 'personal,db,eth,net,web3,admin,miner,txpool' --ipcdisable --port 30302 --rpcaddr 127.0.0.1 --rpcport 8102 --allow-insecure-unlock console 2>> /tmp/eth/60/02.log `

* Next clone this project and build it - `go build` and ensure there are no errors.

* Next just run `go test`
  * A few messages will scroll by and you will see the message "DOUBLE SPEND DETECTED" scroll by. 
  * We have simulated a double spend on our local 2 node cluster

* What happens when you run `go test`
  * In the eth module there is just a single test in the ds_test.go file which gets executed when you run `cd eth && go test`
  * This is what the test does
    * We make node1 be the bad one trying to do malicious things and node2 be the honest one.
    * Create a new address for Malice (the bad one trying the double spend), Bob, Jane and two accounts for getting the mining rewards eb1 (for node1) and eb2 (for node2)
    * Set etherbase on node 1 to Malice's address. Once mining starts, Malice will get the rewards initially. 
    * Set etherbase on node 2 to Bob's address. Bob gets the mining rewards on node2 initially.
    * Add peer of node 2 on node 1 so they sync
    * Start mining for 20 seconds so Malice and Bob have some ETH
    * Stop mining on both nodes
    * Set etherbase accounts for node 1 and node 2 to the eb1 and eb2 respectively. So Malice and Bob won't get the rewards anymore and they have some ETH they already got from the mining earlier. 
    * Start full blockchain scan for Malice's address on the good node - node2
    * Subscribe to new block notifications on node2
    * Remove node 2 peer from node 1. This is so that Malice can play her tricks without anyone in the network noticing.
    * Transfer 50% of the ETH from Malice to Jane and send transaction to node2. Node1 will not see this transaction because node2 is no longer a peer of node1. 
    * Start mining again on node2 and after 10 seconds Malice would have transferred to Jane successfully on node2 while node1 knows nothing about it. Let us assume Jane has shipped the expensive stuff to Malice after the confirmation.
    * We can stop mining on node2 now temporarily again.
    * Next, Malice sends 80% of her original balance (before the transaction with Jane) to node1 only
    * Start mining on node1 and node2. However, Malice will try to overcome the hashpower of the rest of the network. So node1 will mine with far more threads than node2.
    * Then, node2 is added as a peer of node1 again.
    * TADA - Double spend detected
      * Internally, when the full blockchain scan happens for Malice's address, it keeps the history of Malice's transactions done
      * When new block notifications arrive, all cached transactions are run through and the transaction receipt fetched. 
      * If the transaction receipt cannot be fetched anymore, it is deemed a doublespend.
      * The newly arrived transactions are then added to the cache again if they are of interest to the system.

