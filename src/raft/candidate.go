package raft

import (
	"math/rand"
	"time"
)

func (state *ServerState) ToCandidate() {
	oldState := state.state
	state.state = CandidateState
	state.lastElectionTime = time.Now()
	DPrintf("%v: ToCandidate from %v", state, oldState)
}

func (rf *Raft) RunCandidate() {
	rf.RunElection()
	electionTimeout := time.Duration(50+(rand.Int63()%300)) * time.Millisecond
	time.Sleep(electionTimeout)
}
