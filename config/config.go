package config

import (
	"errors"
	"github.com/spf13/viper"
)

type Config struct {
	System    *System    `mapstructure:"system"`
	Showstart *Showstart `mapstructure:"showstart"`
	Ticket    *Ticket    `mapstructure:"ticket"`
}

type System struct {
	MaxGoroutine int `mapstructure:"max_goroutine"`
	MinInterval  int `mapstructure:"min_interval"`
	MaxInterval  int `mapstructure:"max_interval"`
}

type Showstart struct {
	Sign        string `mapstructure:"sign"`
	Token       string `mapstructure:"token"`
	Cookie      string `mapstructure:"cookie"`
	StFlpv      string `mapstructure:"st_flpv"`
	Cusid       string `mapstructure:"cusid"`
	Cusname     string `mapstructure:"cusname"`
	Cversion    string `mapstructure:"cversion"`
	Cterminal   string `mapstructure:"cterminal"`
	Cdeviceinfo string `mapstructure:"cdeviceinfo"`
}

type Ticket struct {
	ActivityId int          `mapstructure:"activity_id"`
	StartTime  string       `mapstructure:"start_time"`
	List       []TicketList `mapstructure:"list"`
	People     []string     `mapstructure:"people"`
}

type TicketList struct {
	Session string `mapstructure:"session"`
	Price   string `mapstructure:"price"`
}

func InitCfg() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./..")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("cant not find config.mapstructure")
		}
	}

	var cfg *Config

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, errors.New("cant not find config.mapstructure")
	}

	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *Config) Validate() error {
	if cfg.Ticket == nil {
		return errors.New("未读取到票务配置信息")
	}

	if len(cfg.Ticket.List) == 0 {
		return errors.New("未读取到要购票的场次以及票价")
	}

	if len(cfg.Ticket.People) == 0 {
		return errors.New("未读取到观演人信息")
	}

	return nil
}
