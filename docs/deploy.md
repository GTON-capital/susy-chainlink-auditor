

```
git clone https://github.com/SuSy-One/susy-chainlink-auditor
cd susy-chainlink-auditor/
mkdir chainlink_data
cd chainlink_data
```

создаем 2 файла

.password - в нем содержится пароль доступа
```
echo "MyVeryStrongP@ssword" > .password
```

.api - параметры доступа к api ноды
```
myemail@email.com
StringPasswordForApi
```
в файле chainlink.env пишем параметры запуска ноды

запускаем ноду 
```
docker-compose up -d
```

и заходим на ноду http://localhost:7788

вводим креды из файла .api

Попадаем в панель управления нодой. В разделе Keys в самом низу берем адрес ноды в чейне и сообщаем владельцу фида, чтобы авторизовать оракл. Также на этот адрес закидываем токены для оплаты транзакций.

Заходим в раздел Jobs, жмем *New Job*, выбираем формат TOML и создаем описание job

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

Вместо DATA_FEED_CONTRACT_ADDRESS пишем адрес контракта фида, возвращаемся в раздел Jobs, выбираем только созданную Job и смотрим ее Definition, 
копируем externalJobID.

Редактируем файл docker-compose.yaml - раздел *feed-initiator* секцию *environment*:

```
- SUSY_FEED_EMAIL=email@email.com #логин из файла .api
- SUSY_FEED_PASSWORD=password #пароль из файла .api
- SUSY_FEED_JOB_ID=244b8cc7-d791-4c98-a307-4b4affda6923 #сюда вставляем  externalJobID
- SUSY_FEED_CHAIN_URL=https://rpc.ftm.tools/
- SUSY_FEED_BLOCKS_FRAME=100 #количество блоков в раунде
- SUSY_FEED_DURATION=2h #время между пушами в фид
```

перезапускаеи инициатор 

```
docker-compose up -d feed-initiator
```
