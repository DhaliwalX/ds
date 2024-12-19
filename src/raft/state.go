package raft

import (
	"fmt"
	"time"
)

// Current state of the server
type raftState int

const (
	CandidateState raftState = iota
	FollowerState
	LeaderState
)

func (state raftState) String() string {
	switch state {
	case CandidateState:
		return "C"
	case FollowerState:
		return "F"
	case LeaderState:
		return "L"
	}
	return "U"
}

// Log entry for the state machine.
type LogEntry struct {
	Term    int
	Command interface{}
}

type ServerState struct {
	numPeers int
	me       int

	lastElectionTime time.Time

	// Persistent state on all servers.
	state       raftState
	currentTerm int

	// Candidate voted in this term. -1 if none.
	votedFor int
	log      []LogEntry

	// Volatile state on all servers.
	commitIndex int
	lastApplied int

	// Volatile state on leaders.
	nextIndex  []int
	matchIndex []int
}

func (state *ServerState) Initialize(numPeers int, me int) {
	state.numPeers = numPeers
	state.me = me
	state.state = FollowerState
	state.currentTerm = 0
	state.votedFor = -1
	state.log = make([]LogEntry, 0)
	state.commitIndex = 0
	state.lastApplied = 0
}

func (state *ServerState) ToFollower(term int) {
	state.state = FollowerState
	state.currentTerm = term
	state.votedFor = -1
	state.lastElectionTime = time.Now()
	DPrintf("%v\n", state)
}

func (state *ServerState) ToLeader() {
	state.state = LeaderState
	state.votedFor = -1

	state.nextIndex = make([]int, state.numPeers)
	state.matchIndex = make([]int, state.numPeers)
	for i := 0; i < state.numPeers; i++ {
		state.nextIndex[i] = len(state.log)
		state.matchIndex[i] = 0
	}
	DPrintf("%v\n", state)
}

func (state *ServerState) ToCandidate() {
	state.state = CandidateState
	state.lastElectionTime = time.Now()
	DPrintf("%v\n", state)
}

func (state *ServerState) IsLogUptoDate(logIndex, logTerm int) bool {
	if len(state.log) == 0 {
		return true
	}

	lastLog := state.log[len(state.log)-1]

	// If the logs end with the same term, the longer log is more up-to-date.
	if lastLog.Term == logTerm {
		return len(state.log) < logIndex
	}

	// Else, the log with the higher term is more up-to-date.
	return lastLog.Term < logTerm
}

func (state *ServerState) LastLogIndex() int {
	return len(state.log) - 1
}

func (state *ServerState) LastLogTerm() int {
	if state.LastLogIndex() < 0 {
		return -1
	}
	return state.log[state.LastLogIndex()].Term
}

func (state *ServerState) Append(entries ...LogEntry) {
	state.log = append(state.log, entries...)
}

func (state *ServerState) String() string {
	str := []byte{}
	str = fmt.Appendf(str, "CurrentTerm: %d, ", state.currentTerm)
	str = fmt.Appendf(str, "Me: %d, ", state.me)
	str = fmt.Appendf(str, "State: %s, ", state.state)
	str = fmt.Appendf(str, "VotedFor: %d, ", state.votedFor)
	str = fmt.Appendf(str, "CommitIndex: %d, ", state.commitIndex)
	str = fmt.Appendf(str, "LastApplied: %d, ", state.lastApplied)
	str = fmt.Appendf(str, "Log: %v, ", state.log)
	return string(str)
}

/*
Apna College courses have provided me with a strong foundation in C++ and Web development. I am now ready to move beyond theoretical knowledge and gain practical experience, which is why I am excited about this internship. I believe this internship will provide the perfect environment for me to further develop my skills, ultimately allowing me to contribute back to the community that has helped me grow.
*/
