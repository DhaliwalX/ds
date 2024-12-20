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

func (entry LogEntry) String() string {
	return fmt.Sprintf("{T: %d C: %v}", entry.Term, entry.Command)
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
	state.commitIndex = 0
	state.lastApplied = 0
	state.log = append(state.log, LogEntry{Term: state.currentTerm, Command: nil})
}

func (state *ServerState) IsLogUptoDate(logIndex, logTerm int) bool {
	lastLog := state.log[state.LastLogIndex()]

	// If the logs end with the same term, the longer log is more up-to-date.
	if lastLog.Term == logTerm {
		return state.LastLogIndex() <= logIndex
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
	if state.state == LeaderState {
		str = fmt.Appendf(str, "NextIndex: %v, ", state.nextIndex)
		str = fmt.Appendf(str, "MatchIndex: %v, ", state.matchIndex)
	}
	return string(str)
}
