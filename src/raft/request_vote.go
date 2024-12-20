package raft

import (
	"fmt"
	"time"
)

// example RequestVote RPC arguments structure.
// field names must start with capital letters!
type RequestVoteArgs struct {
	// Your data here (3A, 3B).
	Term         int
	CandidateId  int
	LastLogIndex int
	LastLogTerm  int
}

// example RequestVote RPC reply structure.
// field names must start with capital letters!
type RequestVoteReply struct {
	// Your data here (3A).
	Term        int
	VoteGranted bool
}

func (reply *RequestVoteReply) Reject(term int) {
	reply.Term = term
	reply.VoteGranted = false
}

func (reply *RequestVoteReply) Grant(term int) {
	reply.Term = term
	reply.VoteGranted = true
}

func (args *RequestVoteArgs) String() string {
	return fmt.Sprintf("RV{Term: %d, CandidateId: %d, LastLogIndex: %d, LastLogTerm: %d}",
		args.Term, args.CandidateId, args.LastLogIndex, args.LastLogTerm)
}

func (args *RequestVoteReply) String() string {
	return fmt.Sprintf("RV{Term: %d, VoteGranted: %v}", args.Term, args.VoteGranted)
}

// example RequestVote RPC handler.
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	// Your code here (3A, 3B).
	rf.mu.Lock()
	defer rf.mu.Unlock()
	defer DPrintf("%v: RequestVote(%v) -> %v", &rf.state, args, reply)
	if args.Term < rf.Term() {
		reply.Reject(rf.Term())
		return
	}

	defer rf.persist()
	if args.Term > rf.Term() {
		rf.state.ToFollower(args.Term)
	}
	rf.state.currentTerm = args.Term
	rf.state.lastElectionTime = time.Now()
	if (rf.DidNotVote() || rf.VotedFor() == args.CandidateId) &&
		(rf.isLogUptoDate(args.LastLogIndex, args.LastLogTerm)) {
		rf.Vote(args.CandidateId)
		reply.Grant(rf.Term())
	} else {
		reply.Reject(rf.Term())
	}
}
