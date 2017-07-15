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

	consulConfig := api.DefaultConfig()
	if consulConfig.Address != "" {
		consulConfig.Address = config.Address
	}

	consulClient, _ = api.NewClient(consulConfig)

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
			consulClient.Session().Destroy(session, nil)
			session, _, _ = consulClient.Session().Create(sessionEnt, nil)
			kv.Session = session
		}

		time.Sleep(config.WatchWaitTime * time.Second)
	}

}

func checkLeadershipStatus(leader Interface, config *Config) {

	for {

		ok, _, _ := consulClient.KV().Acquire(kv, nil)

		if !ok {
			leader.GiveUpLeadership()
			break
		}

		time.Sleep(config.WatchWaitTime * time.Second)

	}

}
