package config

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	System    *System     `mapstructure:"system"`
	Showstart *Showstart  `mapstructure:"showstart"`
	Ticket    *Ticket     `mapstructure:"ticket"`
	SmtpEmail *smtp_email `mapstructure:"smtp_email"`
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

type smtp_email struct {
	Host     string `mapstructure:"host"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	To       string `mapstructure:"email_to"`
	Enable   bool   `mapstructure:"enable"`
}

func InitCfg() (*Config, error) {
	// 获取当前工作目录
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("未获取到当前工作目录")
	}
	log.Println("当前工作目录：", workDir)

	// 获取可执行文件的路径
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}
	log.Println("可执行文件路径：", exePath)

	// 获取可执行文件的目录
	exeDir := filepath.Dir(exePath)

	// 设置 Viper 的配置文件名和类型
	viper.SetConfigName("config") // 配置文件名（不带扩展名）
	viper.SetConfigType("yaml")   // 配置文件类型

	// 首先尝试从可执行文件目录加载配置
	viper.AddConfigPath(exeDir)

	// 如果在可执行文件目录未找到，则尝试从当前工作目录加载配置
	viper.AddConfigPath(workDir)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("未读取到配置文件，请确认config.yaml是否存在")
		}
	}

	var cfg *Config

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, errors.New("配置信息映射失败，请检查配置文件格式是否遵循yaml格式")
	}

	err = cfg.Validate()
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
