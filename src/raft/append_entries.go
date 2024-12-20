package raft

import (
	"fmt"
	"time"
)

type AppendEntriesArgs struct {
	Term         int
	LeaderId     int
	PrevLogIndex int
	PrevLogTerm  int
	Entries      []LogEntry
	LeaderCommit int
}

type AppendEntriesReply struct {
	Term    int
	Success bool
}

func (args *AppendEntriesArgs) String() string {
	return fmt.Sprintf("AE{Term: %d, LeaderId: %d, PrevLogIndex: %d, PrevLogTerm: %d, Entries: %v, LeaderCommit: %d}",
		args.Term, args.LeaderId, args.PrevLogIndex, args.PrevLogTerm, args.Entries, args.LeaderCommit)
}

func (reply *AppendEntriesReply) String() string {
	return fmt.Sprintf("AE{Term: %d, Success: %v}", reply.Term, reply.Success)
}

func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	defer DPrintf("%v, AppendEntries(%v) -> %v", &rf.state, args, reply)
	rf.state.lastElectionTime = time.Now()
	// DPrintf("%d> AppendEntries(%v, %v) for %v", rf.me, args, reply, rf.Term())
	if args.Term < rf.Term() {
		reply.Success = false
		reply.Term = rf.Term()
		return
	}

	if args.Term > rf.Term() {
		rf.state.ToFollower(args.Term)
	}
	reply.Term = rf.Term()
	defer rf.persist()

	// Log mismatch.
	if (args.PrevLogIndex <= rf.state.LastLogIndex() &&
		rf.state.log[args.PrevLogIndex].Term != args.PrevLogTerm) ||
		(args.PrevLogIndex > rf.state.LastLogIndex()) {
		reply.Success = false
		return
	}

	// Truncate conflicting logs
	rf.state.log = rf.state.log[:args.PrevLogIndex+1]

	// At this stage, logs have matched and we can proceed with the append.
	reply.Success = true

	// Append entries.
	if len(args.Entries) > 0 {
		rf.state.Append(args.Entries...)
	}

	if args.LeaderCommit > rf.state.commitIndex {
		rf.state.commitIndex = min(args.LeaderCommit, rf.state.LastLogIndex())
	}

	go rf.Apply()
}
