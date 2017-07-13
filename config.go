package election

import "time"

type Config struct {
	WatchWaitTime time.Duration
	LeaderKey     string
}

func DefaultConfig() *Config {

	return &Config{
		WatchWaitTime: 1,
		LeaderKey:     "services/leader",
	}
}
