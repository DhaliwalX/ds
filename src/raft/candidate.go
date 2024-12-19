package raft

import (
	"math/rand"
	"time"
)

func (rf *Raft) RunCandidate() {
	rf.RunElection()
	electionTimeout := time.Duration(50+(rand.Int63()%300)) * time.Millisecond
	time.Sleep(electionTimeout)
}
