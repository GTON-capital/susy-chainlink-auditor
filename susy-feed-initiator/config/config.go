package config

type Config struct {
	URL         string `mapstructure:"SUSY_FEED_URL" `
	Email       string `mapstructure:"SUSY_FEED_EMAIL" `
	Password    string `mapstructure:"SUSY_FEED_PASSWORD" `
	JobID       string `mapstructure:"SUSY_FEED_JOB_ID"`
	ChainUrl    string `mapstructure:"SUSY_FEED_CHAIN_URL"`
	BlocksFrame uint32 `mapstructure:"SUSY_FEED_BLOCKS_FRAME"`
	Duration    string `mapstructure:"SUSY_FEED_DURATION"`
	Scheduler   string `mapstructure:"SUSY_FEED_SCHEDULER"`
}

var RuntimeConfig Config
