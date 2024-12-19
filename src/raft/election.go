// Define the election process of the raft algorithm.

package raft

func (rf *Raft) RunElection() {
	rf.mu.Lock()
	rf.state.ToCandidate()
	rf.state.currentTerm++
	rf.Vote(rf.me)
	currentTerm := rf.Term()
	candidateId := rf.me
	lastLogIndex := rf.state.LastLogIndex()
	lastLogTerm := rf.state.LastLogTerm()
	rf.mu.Unlock()

	votes := 1 // 1 for self vote.
	for i := 0; i < rf.state.numPeers; i++ {
		if rf.me == i {
			continue
		}
		go func(peer int) {
			args := RequestVoteArgs{
				Term:         currentTerm,
				CandidateId:  candidateId,
				LastLogIndex: lastLogIndex,
				LastLogTerm:  lastLogTerm,
			}
			reply := RequestVoteReply{}
			//
			if rf.sendRequestVote(peer, &args, &reply) {
				rf.mu.Lock()
				defer rf.mu.Unlock()

				// No longer a candidate
				if rf.state.state != CandidateState {
					return
				}

				if reply.Term > currentTerm {
					rf.state.ToFollower(reply.Term)
				} else if reply.VoteGranted {
					votes++
					if votes > rf.state.numPeers/2 {
						rf.state.ToLeader()
						go rf.RunLeader()
					}
				}
			}
		}(i)
	}
}
