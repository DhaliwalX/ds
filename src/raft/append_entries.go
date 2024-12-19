package raft

import "time"

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

func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	rf.state.lastElectionTime = time.Now()
	// DPrintf("%d> AppendEntries(%v, %v) for %v", rf.me, args, reply, rf.Term())
	DPrintf("%v", rf.state.String())
	if args.Term < rf.Term() {
		reply.Success = false
		reply.Term = rf.Term()
		return
	}

	if args.Term > rf.Term() {
		rf.state.ToFollower(args.Term)
	}
	reply.Term = rf.Term()

	// Log mismatch.
	if (args.PrevLogIndex <= rf.state.LastLogIndex() &&
		rf.state.LastLogTerm() != args.PrevLogTerm) ||
		(args.PrevLogIndex > rf.state.LastLogIndex()) {
		reply.Success = false
		return
	}

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
