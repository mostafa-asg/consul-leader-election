# consul-leader-election
An implementation of Consul leader election

## Install
```
go get github.com/mostafa-asg/consul-leader-election
```

## How to use
You must implement election.Interface interface :
```Go
type Leader struct {

}

//this method is called whenever you become a leader
func (l Leader) TakeLeadership(authority *election.LeadershipAuthority){
  // Do whatever a leader must do  
  fmt.Println("I am a leader")	
}

//this method is called whenever you loose leadership
func (l Leader) GiveUpLeadership(){
  fmt.Println("I'm not a leader anymore")
}
```
after that you can call TryTakingLeadership :
```Go
election.TryTakingLeadership( election.DefaultConfig() , leader)
```
it blocks until you become a leader, if you want non-blocking features you can use :
```Go
go election.TryTakingLeadership( election.DefaultConfig() , leader)
```

## Full Example
```Go
package main

import (
	"fmt"
	"time"
	"github.com/mostafa-asg/consul-leader-election"
)

var iAmALeader bool
var auth *election.LeadershipAuthority

type Leader struct {

}

func (l Leader) TakeLeadership(authority *election.LeadershipAuthority){
  fmt.Println("I am a leader")
  auth = authority
  iAmALeader = true
}

func (l Leader) GiveUpLeadership(){
  fmt.Println("I'm not a leader anymore")
  auth = nil
  iAmALeader = false
  //Try to become a leader again
  election.TryTakingLeadership(electionConfig , leader)
}

var leader Leader
var electionConfig *election.Config

func main() {
  electionConfig = &election.Config{
    WatchWaitTime:1, //Time in seconds to check for leadership
    LeaderKey:"service/myService/leader", // The leadership key to create/aquire
  }
  election.TryTakingLeadership( electionConfig , leader)

  time.Sleep(30 * time.Second)

  if iAmALeader {
    //release leadership lock
    auth.Resignation()		
  }
}

```
