package app

import "flag"

type AgentCfg struct {
	Address        string
	ReportInterval int
	PollInterval   int
}

type ServerCfg struct {
	Address string
}

func getAgentConfig() (cfg AgentCfg) {
	flag.StringVar(&cfg.Address, "a", "localhost:8080", "http server address")
	flag.IntVar(&cfg.ReportInterval, "r", 10, "report interval")
	flag.IntVar(&cfg.PollInterval, "p", 2, "poll interval")
	flag.Parse()

	return cfg
}

func getServerConfig() (cfg ServerCfg) {
	flag.StringVar(&cfg.Address, "a", "localhost:8080", "http server address")
	flag.Parse()

	return cfg
}
