
## Initial setup
```
git clone https://github.com/SuSy-One/susy-chainlink-auditor
cd susy-chainlink-auditor/
mkdir chainlink_data
cd chainlink_data
```
## Access settings

Please note that Chainlink node enforces a password policy that requires a password to have at least 3 numbers and 3 symbols.
First, create two config files in `./chainlink_data`:

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

In Chainlink Operator GUI, visit the Keys section and find your node's blockchain address. Share this address with the feed owner to authorize your oracle. In addition, top up this address with the native token to pay for transactions.

## Job setup

Visit the Jobs section in the GUI, click on *New Job*, select TOML as format and enter a job description:

```toml
type = "webhook"
schemaVersion = 1
name = ""
observationSource = """
parse_round_id  [type=jsonparse path="round_id" data="$(jobRun.requestBody)"]
parse_round_data  [type=jsonparse path="round_data" data="$(jobRun.requestBody)"]

encode_tx    [type=ethabiencode
              abi="pushRoundData(uint32 round, int256 value)"
              data="{\\"round\\": $(parse_round_id),\\"value\\": $(parse_round_data)}"]

submit_tx          [type="ethtx"
               to="DATA_FEED_CONTRACT_ADDRESS"
               data="$(encode_tx)"]
parse_round_id -> parse_round_data -> encode_tx -> submit_tx
"""
```

Replace `DATA_FEED_CONTRACT_ADDRESS` with the address of the feed contract.
Return to the Jobs section, select the created *Job*, and view its *Definition* to copy the `externalJobID`. 

Edit `docker-compose.yaml` file, section *feed-initiator*/*environment*:

```
- SUSY_FEED_EMAIL=email@email.com #login from chainlink_data/.api
- SUSY_FEED_PASSWORD=password #password from chainlink_data/.api
- SUSY_FEED_JOB_ID=244b8cc7-d791-4c98-a307-4b4affda6923 #paste externalJobID from Job Definition
- SUSY_FEED_CHAIN_URL=https://rpc.ftm.tools/
- SUSY_FEED_BLOCKS_FRAME=100 #number of blocks in a round
- SUSY_FEED_SCHEDULER=*/1 * * * * #time interval between consecutive feed pushes
```

Relaunch the feed initiator: 

```
docker-compose up -d feed-initiator
```
