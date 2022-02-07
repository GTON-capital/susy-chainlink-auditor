package sources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"susy-feed-initiator/config"
)

type CoinGeckoDataSource struct {
	client http.Client
}

func NewCoinGeckoDataSource() CoinGeckoDataSource {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	return CoinGeckoDataSource{
		client: *httpClient,
	}
}

func (cg CoinGeckoDataSource) GetFloat64(tick string) float64 {
	switch tick {
	case "gton":
		price, err := cg.loadGtonPrice()
		if err != nil {
			return 0
		}
		return price
	default:
		return 0
	}

}

type HeightRequest struct {
	RpcVer string `json:"jsonrpc"`
	Method string `json:"method"`
	ID     int    `json:"id"`
}
type HeightResponse struct {
	RpcVer string `json:"jsonrpc"`
	Result string `json:"result"`
	ID     int    `json:"id"`
}

func hex2int(hexStr string) uint64 {
	// remove 0x suffix if found in the input string
	cleaned := strings.Replace(hexStr, "0x", "", -1)

	// base 16 for hexadecimal
	result, _ := strconv.ParseUint(cleaned, 16, 64)
	return uint64(result)
}

func (cg CoinGeckoDataSource) GetRoundId() uint64 {
	r := HeightRequest{
		"2.0",
		"eth_blockNumber",
		1,
	}
	rb, _ := json.Marshal(r)
	resp, err := cg.client.Post(config.RuntimeConfig.ChainUrl, "application/json", bytes.NewReader(rb))
	if err != nil {
		return 0
	}
	result := HeightResponse{}
	respb, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(respb, &result)
	height := hex2int(result.Result)
	return height / uint64(config.RuntimeConfig.BlocksFrame)
}

func (cg CoinGeckoDataSource) loadGtonPrice() (float64, error) {
	max_attempts := 5
	for i := 0; i < max_attempts; i++ {
		r, err := cg.client.Get("https://api.coingecko.com/api/v3/simple/token_price/fantom?contract_addresses=0xc1be9a4d5d45beeacae296a7bd5fadbfc14602c4&vs_currencies=usd")
		if err != nil {
			time.Sleep(time.Duration(i) * time.Millisecond * (100))
			continue
		}
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			time.Sleep(time.Duration(i) * time.Millisecond * (100))
			continue
		}
		price := map[string]interface{}{}
		json.Unmarshal(b, &price)
		coin, ok := price["0xc1be9a4d5d45beeacae296a7bd5fadbfc14602c4"]
		if !ok {
			time.Sleep(time.Duration(i) * time.Millisecond * (100))
			continue
		}
		coin2usd, ok := coin.(map[string]interface{})["usd"]
		if !ok {
			time.Sleep(time.Duration(i) * time.Millisecond * (100))
			continue
		}
		return coin2usd.(float64), nil
	}
	return 0, fmt.Errorf("coingecko unavailable")

}
