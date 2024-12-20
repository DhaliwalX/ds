package raft

import (
	"math/rand"
	"time"
)

func (state *ServerState) ToFollower(term int) {
	oldState := state.state
	state.state = FollowerState
	state.currentTerm = term
	state.votedFor = -1
	state.lastElectionTime = time.Now()
	DPrintf("%v: ToFollower from %v\n", state, oldState)
}

func (rf *Raft) RunFollower() {
	electionTimeout := time.Duration(50+(rand.Int63()%300)) * time.Millisecond
	time.Sleep(electionTimeout)

	rf.mu.Lock()
	defer rf.mu.Unlock()
	if rf.State() == FollowerState {
		if time.Since(rf.state.lastElectionTime) >= electionTimeout {
			rf.state.ToCandidate()
		}
	}
}
