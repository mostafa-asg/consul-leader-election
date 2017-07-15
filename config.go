package election

import "time"

type Config struct {
	// Address is the address of the Consul server
	Address       string
	WatchWaitTime time.Duration
	LeaderKey     string
}

func DefaultConfig() *Config {

	return &Config{
		Address:       "127.0.0.1",
		WatchWaitTime: 1,
		LeaderKey:     "services/leader",
	}
}
