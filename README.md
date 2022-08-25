# Job description
```toml
type = "webhook"
schemaVersion = 1
name = ""
externalJobID = "244b8cc7-d791-4c98-a307-4b4affda6927"
observationSource = """
parse_uuid  [type=jsonparse path="uuid" data="$(jobRun.requestBody)"]
parse_source_chain  [type=jsonparse path="source_chain" data="$(jobRun.requestBody)"]
parse_destination_chain  [type=jsonparse path="destination_chain" data="$(jobRun.requestBody)"]
parse_sender  [type=jsonparse path="sender" data="$(jobRun.requestBody)"]
parse_receiver  [type=jsonparse path="receiver" data="$(jobRun.requestBody)"]
parse_amount  [type=jsonparse path="amount" data="$(jobRun.requestBody)"]
parse_destination_tx  [type=jsonparse path="destination_tx" data="$(jobRun.requestBody)"]
parse_source_tx  [type=jsonparse path="source_tx" data="$(jobRun.requestBody)"]

encode_tx    [type=ethabiencode
              abi="addSwap(bytes uuid, bytes sender, string source_chain, bytes receiver, string destination_chain, uint256 amount, bytes source_transaction, bytes destination_transaction)"
              data="{\\"uuid\\": $(parse_uuid),\\"source_chain\\": $(parse_source_chain),\\"destination_chain\\": $(parse_destination_chain),\\"sender\\": $(parse_sender),\\"receiver\\": $(parse_receiver),\\"amount\\": $(parse_amount), \\"destination_transaction\\": $(parse_destination_tx),\\"source_transaction\\": $(parse_source_tx)}"]

submit_tx          [type="ethtx"
               to="0x72E19C4bb9B9f4d23be2243BF26fBDB0F1f59746"
               data="$(encode_tx)"]
parse_uuid -> parse_source_chain -> parse_destination_chain -> parse_sender -> parse_receiver -> parse_amount -> parse_destination_tx -> parse_source_tx -> encode_tx -> submit_tx
"""
```

# Automated deploy

```
curl -O https://raw.githubusercontent.com/GTON-capital/susy-chainlink-auditor/main/deploy.sh && sh deploy.sh
```
