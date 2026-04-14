package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

type AgentCfg struct {
	Address        *string `env:"ADDRESS"`
	ReportInterval *int    `env:"REPORT_INTERVAL"`
	PollInterval   *int    `env:"POLL_INTERVAL"`
}

type ServerCfg struct {
	Address         *string `env:"ADDRESS"`
	StoreInterval   *int    `env:"STORE_INTERVAL"`
	FileStoragePath *string `env:"FILE_STORAGE_PATH"`
	Restore         *bool   `env:"RESTORE"`
	DBConn          *string `env:"DATABASE_DSN"`
}

func GetAgentConfig() (cfg AgentCfg, err error) {
	err = env.Parse(&cfg)
	if err != nil {
		return AgentCfg{}, fmt.Errorf("failed to parse agent env: %w", err)
	}

	if cfg.Address == nil {
		cfg.Address = ptr("localhost:8080")
	}
	flag.StringVar(cfg.Address, "a", *cfg.Address, "http server address")

	if cfg.ReportInterval == nil {
		cfg.ReportInterval = ptr(10)
	}
	flag.IntVar(cfg.ReportInterval, "r", *cfg.ReportInterval, "report interval")

	if cfg.PollInterval == nil {
		cfg.PollInterval = ptr(2)
	}
	flag.IntVar(cfg.PollInterval, "p", *cfg.PollInterval, "poll interval")

	flag.Parse()

	return cfg, nil
}

func GetServerConfig() (cfg ServerCfg, err error) {
	err = env.Parse(&cfg)
	if err != nil {
		return ServerCfg{}, fmt.Errorf("failed to parse server env: %w", err)
	}

	if cfg.Address == nil {
		cfg.Address = ptr("localhost:8080")
	}
	flag.StringVar(cfg.Address, "a", *cfg.Address, "http server address")

	if cfg.StoreInterval == nil {
		cfg.StoreInterval = ptr(300)
	}
	flag.IntVar(cfg.StoreInterval, "i", *cfg.StoreInterval, "store interval")

	if cfg.FileStoragePath == nil {
		cfg.FileStoragePath = ptr("filepath.txt")
	}
	flag.StringVar(cfg.FileStoragePath, "f", *cfg.FileStoragePath, "file storage path")

	if cfg.Restore == nil {
		cfg.Restore = ptr(false)
	}
	flag.BoolVar(cfg.Restore, "r", *cfg.Restore, "restore metrics")

	if cfg.DBConn == nil {
		cfg.DBConn = ptr("postgres://postgres:mypassword@localhost:5432/mydb")
	}
	flag.Parse()

	return cfg, nil
}

func ptr[T any](v T) *T {
	return &v
}
