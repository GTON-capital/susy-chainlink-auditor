package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"susy-feed-initiator/config"
	"susy-feed-initiator/sources"

	"github.com/spf13/viper"

	"github.com/go-co-op/gocron"
)

const curDir = "./"

type InfoRequest struct {
	RoundId   uint64 `json:"round_id"`
	RoundData uint64 `json:"round_data"`
}

func SendInfo(round uint64, value uint64) error {
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
	payload, _ := json.Marshal(creds{Email: config.RuntimeConfig.Email, Password: config.RuntimeConfig.Password})
	r, err := client.Post(config.RuntimeConfig.URL+"/sessions", "application/json", bytes.NewReader(payload))
	if err != nil {
		log.Print(err)
		return err
	}
	fmt.Println(r.Cookies())
	uri := fmt.Sprintf("%s/v2/jobs/%s/runs", config.RuntimeConfig.URL, config.RuntimeConfig.JobID)
	payload, _ = json.Marshal(InfoRequest{round, value})
	fmt.Println(string(payload))
	r, err = client.Post(uri, "application/json", bytes.NewReader(payload))
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func Round() {
	s := sources.NewCoinGeckoDataSource()
	price := uint64(math.Round(s.GetFloat64("gton") * 100))
	round := s.GetRoundId()
	err := SendInfo(round, price)
	if err != nil {
		log.Printf("error: %v", err)
	}
	// log.Print("end polling round")
}

func main() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	// err := viper.ReadInConfig()
	// if err != nil {
	// 	log.Printf("unable to decode into struct, %v", err)
	// }

	t := reflect.TypeOf(config.RuntimeConfig)

	// Iterate over all available fields and read the tag value
	for i := 0; i < t.NumField(); i++ {
		// Get the field, returns https://golang.org/pkg/reflect/#StructField
		field := t.Field(i)

		// Get the field tag value
		tag := field.Tag.Get("mapstructure")
		viper.BindEnv(tag)

	}

	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	err := viper.Unmarshal(&config.RuntimeConfig)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
	}
	s := gocron.NewScheduler(time.UTC)
	s.Cron(config.RuntimeConfig.Scheduler).Do(Round)
	s.StartAsync()

	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan // Blocks here until either SIGINT or SIGTERM is received.

}
