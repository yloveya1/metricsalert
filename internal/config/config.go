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
	Key            *string `env:"KEY"`
	RateLimit      *int    `env:"RATE_LIMIT"`
}

type ServerCfg struct {
	Address         *string `env:"ADDRESS"`
	StoreInterval   *int    `env:"STORE_INTERVAL"`
	FileStoragePath *string `env:"FILE_STORAGE_PATH"`
	Restore         *bool   `env:"RESTORE"`
	DBConn          *string `env:"DATABASE_DSN"`
	Key             *string `env:"KEY"`
}

func GetAgentConfig() (cfg AgentCfg, err error) {
	defaultAddress := "localhost:8080"
	defaultReportInterval := 10
	defaultPollInterval := 2
	defaultRateLimit := 5
	defaultKey := ""

	var address, key string
	var reportInterval, pollInterval, rateLimit int

	flag.StringVar(&address, "a", defaultAddress, "http server address")
	flag.IntVar(&reportInterval, "r", defaultReportInterval, "report interval")
	flag.IntVar(&pollInterval, "p", defaultPollInterval, "poll interval")
	flag.StringVar(&key, "k", defaultKey, "hash key")
	flag.IntVar(&rateLimit, "l", defaultRateLimit, "rate limit")
	flag.Parse()

	err = env.Parse(&cfg)
	if err != nil {
		return AgentCfg{}, fmt.Errorf("failed to parse agent env: %w", err)
	}

	if cfg.Address == nil {
		cfg.Address = &address
	}
	if cfg.ReportInterval == nil {
		cfg.ReportInterval = &reportInterval
	}
	if cfg.PollInterval == nil {
		cfg.PollInterval = &pollInterval
	}
	if cfg.Key == nil {
		cfg.Key = &key
	}
	if cfg.RateLimit == nil {
		cfg.RateLimit = &rateLimit
	}

	return cfg, nil
}

func GetServerConfig() (cfg ServerCfg, err error) {
	defaultAddress := "localhost:8080"
	defaultStoreInterval := 300
	defaultFileStoragePath := "filepath.txt"
	defaultRestore := false
	defaultDBConn := ""
	defaultKey := ""

	var address, fileStoragePath, dbConn, key string
	var storeInterval int
	var restore bool

	flag.StringVar(&address, "a", defaultAddress, "http server address")
	flag.IntVar(&storeInterval, "i", defaultStoreInterval, "store interval")
	flag.StringVar(&fileStoragePath, "f", defaultFileStoragePath, "file storage path")
	flag.BoolVar(&restore, "r", defaultRestore, "restore metrics")
	flag.StringVar(&dbConn, "d", defaultDBConn, "db connection")
	flag.StringVar(&key, "k", defaultKey, "hash key")
	flag.Parse()

	err = env.Parse(&cfg)
	if err != nil {
		return ServerCfg{}, fmt.Errorf("failed to parse server env: %w", err)
	}

	if cfg.Address == nil {
		cfg.Address = &address
	}
	if cfg.StoreInterval == nil {
		cfg.StoreInterval = &storeInterval
	}
	if cfg.FileStoragePath == nil {
		cfg.FileStoragePath = &fileStoragePath
	}
	if cfg.Restore == nil {
		cfg.Restore = &restore
	}
	if cfg.DBConn == nil {
		cfg.DBConn = &dbConn
	}
	if cfg.Key == nil {
		cfg.Key = &key
	}

	return cfg, nil
}
