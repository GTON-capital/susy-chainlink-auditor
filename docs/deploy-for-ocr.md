
## Initial setup
```
git clone https://github.com/SuSy-One/susy-chainlink-auditor
cd susy-chainlink-auditor/
mkdir chainlink-1
cd chainlink-1
```
## Access settings

Please note that Chainlink node enforces a password policy that requires a password to have at least 3 numbers and 3 symbols.
First, create two config files in `./chainlink-1`:

`.password` -- node access password
```
echo "MyVeryStrongP@ssword123" > .password
```

`.api` - login/password for Chainlink Operator GUI access
```
myemail@email.com
StringPasswordForApi123
```

## Node launch
Edit `chainlink.env` and provide node launch parameters. 

Launch the node:
```
docker-compose up -d
```

Use http://localhost:7788 to access Chainlink Operator GUI (or use ssh tunneling if hosting on a remote machine).

Enter the credentials from `.api` file in the GUI.

## Node setup

In Chainlink Operator GUI, visit the Keys section and this node info:

From section Off-Chain Reporting Keys
1. Key ID (looks like 17eb5381e82a6dddc0b19bfbf8abd38538cc1119c3a68208c71b31dd42decd80)
2. Signing Address (looks like ocrsad_0xBca902E8349e61b621e5DcA09AACce4212f91DC5)
3. Peer ID (looks like p2p_12D3KooWJUknXTJfuyFKa7UnJYMSJB811s5ymvBjSQcAjpnSjyBT)
From section EVM Chain Accounts
1. EVM address
2. Operator Wallet address for rewards

Share this info with the feed owner to authorize your oracle. In addition, top up this address with the native token to pay for transactions.

