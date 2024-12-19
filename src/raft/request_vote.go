package raft

import "time"

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

// example RequestVote RPC handler.
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	// Your code here (3A, 3B).
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if args.Term < rf.Term() {
		reply.Reject(rf.Term())
		return
	}

	if args.Term > rf.Term() {
		rf.state.ToFollower(args.Term)
	}
	rf.state.currentTerm = args.Term
	rf.state.lastElectionTime = time.Now()
	if (rf.DidNotVote() || rf.VotedFor() == args.CandidateId) &&
		(rf.isLogUptoDate(args.LastLogIndex, args.LastLogTerm)) {
		rf.Vote(args.CandidateId)
		reply.Grant(rf.Term())

		DPrintf("%v granted vote to %d", &rf.state, args.CandidateId)
	} else {
		reply.Reject(rf.Term())
		DPrintf("%v rejected vote to %d", &rf.state, args.CandidateId)
	}
}
