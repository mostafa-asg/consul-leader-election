package election

import (
	"github.com/hashicorp/consul/api"
	"time"
)

var consulClient *api.Client
var kv *api.KVPair
var session string

type Interface interface {
	TakeLeadership(authority *LeadershipAuthority)
	GiveUpLeadership()
}

func TryTakingLeadership(config *Config, leader Interface) {

	consulClient, _ = api.NewClient(api.DefaultConfig())

	sessionEnt := &api.SessionEntry{
		Name:     "masterElectionSession",
		Behavior: api.SessionBehaviorDelete,
	}

	session, _, _ = consulClient.Session().Create(sessionEnt, nil)

	agent, _ := consulClient.Agent().Self()

	kv = &api.KVPair{
		Key:     config.LeaderKey,
		Value:   []byte(agent["Config"]["NodeName"].(string)),
		Session: session,
	}

	for {
		ok, _, _ := consulClient.KV().Acquire(kv, nil)
		if ok {
			leader.TakeLeadership(new(LeadershipAuthority))
			go checkLeadershipStatus(leader, config)
			break
		} else {
			time.Sleep(config.WatchWaitTime * time.Second)
		}
	}

}

func checkLeadershipStatus(leader Interface, config *Config) {

	for {

		ok, _, _ := consulClient.KV().Acquire(kv, nil)

		if !ok {
			go leader.GiveUpLeadership()
			break
		}

		time.Sleep(config.WatchWaitTime * time.Second)

	}

}
