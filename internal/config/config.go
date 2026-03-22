package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

type AgentCfg struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

type ServerCfg struct {
	Address string `env:"ADDRESS"`
}

func GetAgentConfig() (cfg AgentCfg, err error) {
	err = env.Parse(&cfg)
	if err != nil {
		return AgentCfg{}, fmt.Errorf("failed to parse agent env: %w", err)
	}

	if cfg.Address == "" {
		flag.StringVar(&cfg.Address, "a", "localhost:8080", "http server address")
	}

	if cfg.ReportInterval == 0 {
		flag.IntVar(&cfg.ReportInterval, "r", 10, "report interval")
	}

	if cfg.PollInterval == 0 {
		flag.IntVar(&cfg.PollInterval, "p", 2, "poll interval")
	}

	fmt.Println(cfg)

	flag.Parse()

	return cfg, nil
}

func GetServerConfig() (cfg ServerCfg, err error) {
	err = env.Parse(&cfg)
	if err != nil {
		return ServerCfg{}, fmt.Errorf("failed to parse server env: %w", err)
	}

	if cfg.Address == "" {
		flag.StringVar(&cfg.Address, "a", "localhost:8080", "http server address")
	}

	flag.Parse()

	return cfg, nil
}
