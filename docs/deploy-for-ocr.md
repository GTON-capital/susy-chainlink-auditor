
## Initial setup
```
git clone https://github.com/SuSy-One/susy-chainlink-auditor
cd susy-chainlink-auditor/
mkdir chainlink-1
cd chainlink-1
```
## Access settings

Please note that Chainlink node enforces a password policy that requires a password to have at least 3 numbers and 3 symbols.
First, create two config files in folder `./chainlink-1`:

`.password` -- node access password (of your choosing):
```
echo "MyVeryStrongP@ssword123" > .password
```

`.api` -- login/password for Chainlink Operator GUI access (both of your choosing, used only locally):
```
touch .api // Crating the file
nano .api // Opening editor
```
Add these two lines in the editor:
```
myemail@email.com
StringPasswordForApi123
```
Press ctrl+X to exit, when prompt asks - save the changes.

## Node launch
Makse sure `chainlink.env` has correct launch parameters - currently there is no need to change anything, repo contains all the setup parameters we need. 

If you are using brand new server - install latest `Docker-compose`, there might be issues with setup using apt-get, better use curl as per [instructions here](https://docs.docker.com/compose/install/).  

Launch the node:
```
docker-compose up -d
```

Check the status of the container using command `docker ps`. There should be several containers containers:
```
CONTAINER ID   IMAGE                                    COMMAND                  CREATED          STATUS                    PORTS                                       NAMES
a46ba5a1f652   susy-chainlink-auditor_node1             "chainlink local nod…"   10 minutes ago   Up 10 minutes (healthy)   0.0.0.0:7788->6688/tcp, :::7788->6688/tcp
2261359f83bd   postgres:13.1-alpine                     "docker-entrypoint.s…"   10 minutes ago   Up 10 minutes             0.0.0.0:5433->5432/tcp, :::5433->5432/tcp
3545acdd1fe0   susy-chainlink-auditor_bridge-peg-usd    "/initiator pegUsd"      10 minutes ago   Up 10 minutes                                           
1d487355a08b   susy-chainlink-auditor_bridge-peg-base   "/initiator pegBase"     10 minutes ago   Up 10 minutes
```
If the status isn't **Up** - but instead something like **Restarting (1)** - it means that there are some configuration problems, you can check the logs: 
```
docker logs susy-chainlink-auditor_node1_1 // enter the name of your failing image here
```
Most likely the issue is with the credentials folder - make sure the naming is correct (as seen above) and the files `.password` and `.api` exist.

Use http://localhost:7788 to access Chainlink Operator GUI if you are running the node on a local machine. If not - use ssh tunneling to be able to launch console using web browser:
```
ssh -L 7788:localhost:7788 root@your.server.ip.address
```
After that you'll be able to open http://localhost:7788 in any browser.

Enter the credentials from `.api` file in the GUI - and you are good to go. If for some reason you need to disable the node you can 

## Node setup

In Chainlink Operator GUI, visit the Keys section (hidden behind the gear) and copy this node info:

From section Off-Chain Reporting Keys
1. Key ID (looks like 17eb5381e82a6dddc0b19bfbf8abd38538cc1119c3a68208c71b31dd42decd80)
2. Signing Address (looks like ocrsad_0xBca902E8349e61b621e5DcA09AACce4212f91DC5)
3. Peer ID (looks like p2p_12D3KooWJUknXTJfuyFKa7UnJYMSJB811s5ymvBjSQcAjpnSjyBT) 

From section EVM Chain Accounts
1. EVM address

Additionally provide operator wallet address for rewards - it shouldn't be set anywhere in the configuration, just provide it to DON operator.
Share this info with the feed owner to authorize your oracle. In addition, top up this address with the native token to pay for transactions.

## Extras
Optionally for increased security install a firewall, on Ubuntu you can use **UFW**
```
sudo apt update
sudo apt upgrade
sudo apt install ufw // in case its not installed
```
You have to make sure UFW allows ssh connections so that you can access the server when you turn on the firewall:
```
sudo ufw app list
```
If the output has **OpenSSH** then you are good to go, otherwise run:
```
sudo ufw allow OpenSSH
```
Then we need to open ports 7788 and 5433 for incoming connections (needed for node operations):
```
sudo ufw allow 7788
sudo ufw allow 5433
```
I've found it useful to reboot the server prior to enabling firewall since for some reason UFW might glitch and keep you out of the server - in this case you might need to use recovery console your hosting provides to disable UFW and be able to ssh into the server.  
When you make sure OpenSSH is in the list of allowed applications - you can turn on the firewall:
```
sudo ufw enable
sudo ufw status // Just to make sure ports are on the list of allowed connections
```
Now all other connections except for the ones allowed by UFW would be denied.
