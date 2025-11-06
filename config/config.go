package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Log struct {
		Level string `yaml:"level"`
		Dir   string `yaml:"dir"`
	} `yaml:"log"`

	Database struct {
		Path string `yaml:"path"`
	} `yaml:"database"`

	Shops []string `yaml:"shops"`

	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`

	Cron struct {
		MachineTypesInterval   int `yaml:"machine_types_interval"`
		MachinesInterval       int `yaml:"machines_interval"`
		MachineDetailsInterval int `yaml:"machine_details_interval"`
	} `yaml:"cron"`
}

var cfg *Config

// Load 加载配置文件
func Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	cfg = &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return err
	}

	return nil
}

// Get 获取配置实例
func Get() *Config {
	return cfg
}

// GetMachineTypesInterval 获取机器类型更新间隔
func GetMachineTypesInterval() time.Duration {
	return time.Duration(cfg.Cron.MachineTypesInterval) * time.Second
}

// GetMachinesInterval 获取机器列表更新间隔
func GetMachinesInterval() time.Duration {
	return time.Duration(cfg.Cron.MachinesInterval) * time.Second
}

// GetMachineDetailsInterval 获取机器详情更新间隔
func GetMachineDetailsInterval() time.Duration {
	return time.Duration(cfg.Cron.MachineDetailsInterval) * time.Second
}
