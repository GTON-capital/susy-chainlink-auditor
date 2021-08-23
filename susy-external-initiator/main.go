package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const curDir = "./"
const RelayABI = "[{\"inputs\":[{\"internalType\":\"contract IWETH\",\"name\":\"_wnative\",\"type\":\"address\"},{\"internalType\":\"contract IUniswapV2Router01\",\"name\":\"_router\",\"type\":\"address\"},{\"internalType\":\"contract IERC20\",\"name\":\"_gton\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_relayTopic\",\"type\":\"bytes32\"},{\"internalType\":\"string[]\",\"name\":\"allowedChains\",\"type\":\"string[]\"},{\"internalType\":\"uint256[2][]\",\"name\":\"fees\",\"type\":\"uint256[2][]\"},{\"internalType\":\"uint256[2][]\",\"name\":\"limits\",\"type\":\"uint256[2][]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"feeMin\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"feePercent\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountMinusFee\",\"type\":\"uint256\"}],\"name\":\"CalculateFee\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"}],\"name\":\"DeliverRelay\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"destinationHash\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"bytes\",\"name\":\"receiverHash\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"destination\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Lock\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"uuid\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"chain\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"emiter\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"bytes\",\"name\":\"token\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"RouteValue\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"parser\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"newBool\",\"type\":\"bool\"}],\"name\":\"SetCanRoute\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"destination\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_feeMin\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_feePercent\",\"type\":\"uint256\"}],\"name\":\"SetFees\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"chain\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"newBool\",\"type\":\"bool\"}],\"name\":\"SetIsAllowedChain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"destination\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_lowerLimit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_upperLimit\",\"type\":\"uint256\"}],\"name\":\"SetLimits\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"ownerOld\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"ownerNew\",\"type\":\"address\"}],\"name\":\"SetOwner\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"topicOld\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"topicNew\",\"type\":\"bytes32\"}],\"name\":\"SetRelayTopic\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"walletOld\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"walletNew\",\"type\":\"address\"}],\"name\":\"SetWallet\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"canRoute\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"feeMin\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"feePercent\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gton\",\"outputs\":[{\"internalType\":\"contract IERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"isAllowedChain\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"destination\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"}],\"name\":\"lock\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"lowerLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contract IERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"reclaimERC20\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"reclaimNative\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"relayTopic\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"uuid\",\"type\":\"bytes16\"},{\"internalType\":\"string\",\"name\":\"chain\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"emiter\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"token\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"routeValue\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"router\",\"outputs\":[{\"internalType\":\"contract IUniswapV2Router01\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"parser\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"_canRoute\",\"type\":\"bool\"}],\"name\":\"setCanRoute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"destination\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_feeMin\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_feePercent\",\"type\":\"uint256\"}],\"name\":\"setFees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"chain\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"newBool\",\"type\":\"bool\"}],\"name\":\"setIsAllowedChain\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"destination\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"_lowerLimit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_upperLimit\",\"type\":\"uint256\"}],\"name\":\"setLimits\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"setOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_relayTopic\",\"type\":\"bytes32\"}],\"name\":\"setRelayTopic\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"upperLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"wnative\",\"outputs\":[{\"internalType\":\"contract IWETH\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]"
const RelayParserABI = "[{\"inputs\":[{\"internalType\":\"contract IOracleRouterV2\",\"name\":\"_router\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_nebula\",\"type\":\"address\"},{\"internalType\":\"string[]\",\"name\":\"evmChains\",\"type\":\"string[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"nebula\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes16\",\"name\":\"uuid\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"chain\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"emiter\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"topic0\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"token\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"sender\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"receiver\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"AttachValue\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"chain\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"newBool\",\"type\":\"bool\"}],\"name\":\"SetIsEVM\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"nebulaOld\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"nebulaNew\",\"type\":\"address\"}],\"name\":\"SetNebula\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"ownerOld\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"ownerNew\",\"type\":\"address\"}],\"name\":\"SetOwner\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"contract IOracleRouterV2\",\"name\":\"routerOld\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contract IOracleRouterV2\",\"name\":\"routerNew\",\"type\":\"address\"}],\"name\":\"SetRouter\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"attachValue\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"b\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"offset\",\"type\":\"uint256\"}],\"name\":\"bytesToBytes16\",\"outputs\":[{\"internalType\":\"bytes16\",\"name\":\"\",\"type\":\"bytes16\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"b\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"offset\",\"type\":\"uint256\"}],\"name\":\"bytesToBytes32\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"b\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"startPos\",\"type\":\"uint256\"}],\"name\":\"deserializeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"b\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"startPos\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"len\",\"type\":\"uint256\"}],\"name\":\"deserializeUint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"a\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"b\",\"type\":\"string\"}],\"name\":\"equal\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"isEVM\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nebula\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"router\",\"outputs\":[{\"internalType\":\"contract IOracleRouterV2\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"chain\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"newBool\",\"type\":\"bool\"}],\"name\":\"setIsEVM\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_nebula\",\"type\":\"address\"}],\"name\":\"setNebula\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"setOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contract IOracleRouterV2\",\"name\":\"_router\",\"type\":\"address\"}],\"name\":\"setRouter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"\",\"type\":\"bytes16\"}],\"name\":\"uuidIsProcessed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

type Lock struct {
	DestinationChain string         //:PLG
	Receiver         common.Address //:F87A9819CE260FB710C00BB841BF4B8B311EC741
	Amount           uint64         //:461082
	TxHash           common.Hash
}

type AttachValue struct {
	UUID        [16]byte
	SourceChain string         //:PLG
	Sender      common.Address //:CED486E3905F8FE1E8AF5D1791F5E7AD7915F01A000000000000000000000000
	Receiver    common.Address //:CED486E3905F8FE1E8AF5D1791F5E7AD7915F01A000000000000000000000000
	Amount      uint64         //:299465573416930282
	TxHash      common.Hash
}

type RouteInfo struct {
	UUID             [16]byte
	DestinationChain string
	Amount           uint64 //:299465573416930282
	TxHash           common.Hash
}

// type SwapInfoSender struct {
// 	client *http.Client
// }

type SwapInfo struct {
	UUID              [16]byte       `json:"uuid"`
	SourceChain       string         `json:"source_chain"` //:PLG
	DestinationChain  string         `json:"destination_chain"`
	Sender            common.Address `json:"sender"`   //:CED486E3905F8FE1E8AF5D1791F5E7AD7915F01A000000000000000000000000
	Receiver          common.Address `json:"receiver"` //:CED486E3905F8FE1E8AF5D1791F5E7AD7915F01A000000000000000000000000
	Amount            uint64         `json:"amount"`   //:299465573416930282
	DestinationTxHash common.Hash    `json:"destination_tx"`
	SourceTxHash      common.Hash    `json:"source_tx"`
}

func (s SwapInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		UUID              string         `json:"uuid"`
		SourceChain       string         `json:"source_chain"` //:PLG
		DestinationChain  string         `json:"destination_chain"`
		Sender            common.Address `json:"sender"`   //:CED486E3905F8FE1E8AF5D1791F5E7AD7915F01A000000000000000000000000
		Receiver          common.Address `json:"receiver"` //:CED486E3905F8FE1E8AF5D1791F5E7AD7915F01A000000000000000000000000
		Amount            uint64         `json:"amount"`   //:299465573416930282
		DestinationTxHash common.Hash    `json:"destination_tx"`
		SourceTxHash      common.Hash    `json:"source_tx"`
	}{
		UUID:              "0x" + hex.EncodeToString(s.UUID[:]),
		SourceChain:       s.SourceChain,
		DestinationChain:  s.DestinationChain,
		Sender:            s.Sender,
		Receiver:          s.Receiver,
		Amount:            s.Amount,
		DestinationTxHash: s.DestinationTxHash,
		SourceTxHash:      s.SourceTxHash,
	})
}

// type LockVisited struct {
// 	Chain            string `mapstructure:"chainname"`
// 	LastVisitedBlock uint64 `mapstructure:"lastvisitedblock"`
// }

type ChainConfig struct {
	ChainName        string `mapstructure:"chainname"`
	Endpoint         string `mapstructure:"endpoint"`
	AttMu            sync.Mutex
	LastVisitedBlock uint64 `mapstructure:"lastvisitedblock"`
	LockDepth        uint64 `mapstructure:"lockdepth"`
	FrameSize        uint64 `mapstructure:"framesize"`
	RelayContract    string `mapstructure:"relaycontract"`
	ParserContract   string `mapstructure:"parsercontract"`
}

type Config struct {
	Api           ApiConfig               `mapstructure:"api"`
	Chains        map[string]*ChainConfig `mapstructure:"chains"`
	MonitorChains []string                `mapstructure:"monitorchains"`
}

type ApiConfig struct {
	URL      string `mapstructure:"url"`
	Email    string `mapstructure:"email"`
	Password string `mapstructure:"password"`
	JobID    string `mapstructure:"jobid"`
}

var RuntimeConfig Config

func SendSwapInfo(sw SwapInfo) error {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Print(err)
		return err
	} // error handling }

	client := &http.Client{
		Jar: jar,
	}

	type creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	payload, _ := json.Marshal(creds{Email: RuntimeConfig.Api.Email, Password: RuntimeConfig.Api.Password})
	r, err := client.Post(RuntimeConfig.Api.URL+"/sessions", "application/json", bytes.NewReader(payload))
	if err != nil {
		log.Print(err)
		return err
	}
	fmt.Println(r.Cookies())
	uri := fmt.Sprintf("%s/v2/jobs/%s/runs", RuntimeConfig.Api.URL, RuntimeConfig.Api.JobID)
	payload, _ = json.Marshal(sw)
	r, err = client.Post(uri, "application/json", bytes.NewReader(payload))
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

// func (s *SwapInfoSender) Initizlize() {
// 	jar, err := cookiejar.New(nil)
// 	if err != nil {
// 		log.Print(err)
// 	} // error handling }

// 	s.client = &http.Client{
// 		Jar: jar,
// 	}

// 	type creds struct {
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 	}
// 	payload, _ := json.Marshal(creds{Email: RuntimeConfig.Api.Email, Password: RuntimeConfig.Api.Password})
// 	r, err := s.client.Post(RuntimeConfig.Api.URL+"/sessions", "application/json", bytes.NewReader(payload))
// 	if err != nil {
// 		log.Print(err)
// 	}
// 	fmt.Println(r.Cookies())
// 	uri := fmt.Sprintf("%s/v2/jobs/%s/runs", RuntimeConfig.Api.URL, RuntimeConfig.Api.JobID)
// 	r, err = s.client.Post(uri, "application/json", bytes.NewReader([]byte("500")))
// 	if err != nil {
// 		log.Print(err)
// 	}
// 	fmt.Println(r.Cookies())
// }

func GetAttaches(client *ethclient.Client, address string, startBlock int64, framesize int64) ([]AttachValue, uint64, error) {
	res := []AttachValue{}

	var lastBlock uint64
	contractAbi, err := abi.JSON(strings.NewReader(string(RelayParserABI)))
	if err != nil {
		return []AttachValue{}, 0, err
	}

	lastBlock, err = client.BlockNumber(context.Background())
	if err != nil {
		return []AttachValue{}, 0, err
	}
	//var wg sync.WaitGroup
	for i := startBlock; i < int64(lastBlock); i += framesize {
		contractAddress := common.HexToAddress(address)
		endBlock := uint64(i + framesize)
		if endBlock > lastBlock {
			endBlock = lastBlock
		}
		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(i)),
			ToBlock:   big.NewInt(int64(endBlock)),
			Addresses: []common.Address{
				contractAddress,
			},
			Topics: [][]common.Hash{
				{
					common.HexToHash("0xdedc6edbc70712376c85a9aafdc87bbce69a70b6f789a345558447c88785b553"),
				},
			},
		}

		logs, err := client.FilterLogs(context.Background(), query)
		if err != nil {
			return []AttachValue{}, 0, err
		}
		_res := []AttachValue{}
		for _, vLog := range logs {

			event, err := contractAbi.Unpack("AttachValue", vLog.Data)
			if err != nil {
				return []AttachValue{}, 0, err
			}
			attVal := AttachValue{
				UUID:        event[1].([16]uint8),
				SourceChain: event[2].(string),
				Sender:      common.BytesToAddress(event[6].([]uint8)[:20]),
				Receiver:    common.BytesToAddress(event[7].([]uint8)[:20]),
				Amount:      event[8].(*big.Int).Uint64(),
				TxHash:      vLog.TxHash,
			}
			_res = append(_res, attVal)
		}
		res = append(res, _res...)
	}
	return res, lastBlock, nil
}

func GetRoutes(client *ethclient.Client, address string, startBlock int64, framesize int64) ([]RouteInfo, uint64, error) {
	res := []RouteInfo{}

	var lastBlock uint64
	contractAbi, err := abi.JSON(strings.NewReader(string(RelayABI)))
	if err != nil {
		return []RouteInfo{}, 0, err
	}

	lastBlock, err = client.BlockNumber(context.Background())
	if err != nil {
		return []RouteInfo{}, 0, err
	}
	for i := startBlock; i < int64(lastBlock); i += framesize {
		contractAddress := common.HexToAddress(address)

		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(i)),
			ToBlock:   big.NewInt(int64(i + int64(framesize))),
			Addresses: []common.Address{
				contractAddress,
			},
			Topics: [][]common.Hash{
				{
					common.HexToHash("0x112b98d9b7f0fd96e462b22deffd3ec2c95405e8458b249860f64e8a6ebf4b59"),
				},
			},
		}

		logs, err := client.FilterLogs(context.Background(), query)
		if err != nil {
			return []RouteInfo{}, 0, err
		}
		_res := []RouteInfo{}
		for _, vLog := range logs {

			event, err := contractAbi.Unpack("RouteValue", vLog.Data)
			if err != nil {
				return []RouteInfo{}, 0, err
			}
			dc := event[1].(string)
			dc = string(bytes.Trim([]byte(dc), "\x00"))
			attVal := RouteInfo{
				UUID:             event[0].([16]uint8),
				DestinationChain: dc,
				Amount:           event[3].(*big.Int).Uint64(),
				TxHash:           vLog.TxHash,
			}
			_res = append(_res, attVal)
		}
		res = append(res, _res...)
	}
	return res, lastBlock, nil
}

func GetLocks(client *ethclient.Client, address string, startBlock int64, framesize int64) ([]Lock, uint64, error) {
	res := []Lock{}

	var lastBlock uint64
	contractAbi, err := abi.JSON(strings.NewReader(string(RelayABI)))
	if err != nil {
		return []Lock{}, 0, err
	}

	lastBlock, err = client.BlockNumber(context.Background())
	if err != nil {
		return []Lock{}, 0, err
	}
	for i := startBlock; i < int64(lastBlock); i += framesize {
		contractAddress := common.HexToAddress(address)

		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(i)),
			ToBlock:   big.NewInt(int64(i + framesize)),
			Addresses: []common.Address{
				contractAddress,
			},
			Topics: [][]common.Hash{
				{
					common.HexToHash("0xa4f88aed847e87bafdc18210d88464dc24f71fa4bf1b4672710c9bc876bb0044"),
				},
			},
		}

		logs, err := client.FilterLogs(context.Background(), query)
		if err != nil {
			return []Lock{}, 0, err
		}
		_res := []Lock{}
		for _, vLog := range logs {

			event, err := contractAbi.Unpack("Lock", vLog.Data)
			if err != nil {
				return []Lock{}, 0, err
			}

			attVal := Lock{
				DestinationChain: event[0].(string),
				Receiver:         common.BytesToAddress(event[1].([]uint8)[:20]),
				Amount:           event[2].(*big.Int).Uint64(),
				TxHash:           vLog.TxHash,
			}
			_res = append(_res, attVal)
		}
		res = append(res, _res...)
	}
	return res, lastBlock, nil
}

func MatchSwaps(vals []AttachValue, routes []RouteInfo, locks []Lock) []SwapInfo {
	res := []SwapInfo{}

	attMap := make(map[[16]uint8]AttachValue)
	rtMap := make(map[[16]uint8]RouteInfo)
	for _, av := range vals {
		attMap[av.UUID] = av
	}
	for _, rt := range routes {
		rtMap[rt.UUID] = rt
	}

	for id, val := range attMap {
		fmt.Printf("Attach: %s", val.TxHash.String())
		rt, ok := rtMap[id]
		if !ok {
			continue
		}

		found := false
		lock := Lock{}
		for _, lk := range locks {
			if lk.Amount == rt.Amount &&
				lk.DestinationChain == rt.DestinationChain &&
				lk.Receiver == val.Receiver {
				lock = lk
				found = true
			}
		}
		if !found {
			continue
		}
		res = append(res, SwapInfo{
			UUID:              val.UUID,
			SourceChain:       val.SourceChain,
			DestinationChain:  rt.DestinationChain,
			Sender:            val.Sender,
			Receiver:          val.Receiver,
			Amount:            val.Amount,
			DestinationTxHash: rt.TxHash,
			SourceTxHash:      lock.TxHash,
		})
	}

	return res
}

func Round() {
	log.Println("Start polling chains...")
	for _, chain := range RuntimeConfig.MonitorChains {
		log.Printf("\t\tPolling %s", chain)
		locksCache := make(map[string][]Lock)

		cConf, ok := RuntimeConfig.Chains[chain]
		if !ok {
			continue
		}
		cConf.AttMu.Lock()

		client, err := ethclient.Dial(cConf.Endpoint)
		if err != nil {
			cConf.AttMu.Unlock()
			continue
		}
		attaches, lastBLock, err := GetAttaches(client, cConf.ParserContract, int64(cConf.LastVisitedBlock), int64(cConf.FrameSize))
		if err != nil {
			cConf.AttMu.Unlock()
			continue
		}
		log.Printf("\t\t%s - found %d attaches", chain, len(attaches))
		routes, _, err := GetRoutes(client, cConf.RelayContract, int64(cConf.LastVisitedBlock), int64(cConf.FrameSize))
		if err != nil {
			cConf.AttMu.Unlock()
			continue
		}

		for _, attVal := range attaches {
			source_chain := strings.ToLower(attVal.SourceChain)
			locks, ok := locksCache[source_chain]
			if !ok {
				locksConfig, ok := RuntimeConfig.Chains[source_chain]
				if !ok {
					continue
				}
				locksClient, err := ethclient.Dial(locksConfig.Endpoint)
				if err != nil {
					continue
				}
				last, err := locksClient.BlockNumber(context.Background())
				if err != nil {
					continue
				}
				_locks, _, err := GetLocks(locksClient, locksConfig.RelayContract, int64(last-locksConfig.LockDepth), int64(locksConfig.FrameSize))
				if err != nil {
					continue
				}
				locks = _locks
				locksCache[source_chain] = locks
			}

			swaps := MatchSwaps(attaches, routes, locks)
			for _, sw := range swaps {
				err := SendSwapInfo(sw)
				if err != nil {
					log.Println(err)
				}
			}
		}
		cConf.LastVisitedBlock = lastBLock
		bs, _ := yaml.Marshal(RuntimeConfig)
		viper.MergeConfig(bytes.NewReader(bs))

		viper.WriteConfig()
		cConf.AttMu.Unlock()

	}
	log.Print("end polling round")
}

func main() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	err = viper.Unmarshal(&RuntimeConfig)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	go func() {
		for {
			Round()
			time.Sleep(time.Minute * 5)
		}
	}()

	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan // Blocks here until either SIGINT or SIGTERM is received.

}
