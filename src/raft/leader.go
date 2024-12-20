package raft

import "time"

func (state *ServerState) ToLeader() {
	oldState := state.state
	state.state = LeaderState
	state.votedFor = -1

	state.nextIndex = make([]int, state.numPeers)
	state.matchIndex = make([]int, state.numPeers)
	for i := 0; i < state.numPeers; i++ {
		state.nextIndex[i] = len(state.log)
		state.matchIndex[i] = 0
	}
	DPrintf("%v: ToLeader from %v", state, oldState)
}

func (rf *Raft) safeTerm(index int) int {
	if index < 0 {
		return -1
	}
	return rf.state.log[index].Term
}

func (rf *Raft) safeEntries(index int) []LogEntry {
	if index < 0 {
		index = 0
	}
	return rf.state.log[index:]
}

func (rf *Raft) BroadcastHeartbeat() {
	for i := 0; i < rf.state.numPeers; i++ {
		if i == rf.me {
			continue
		}

		go func(peer int) {
			rf.mu.Lock()
			if rf.state.state != LeaderState {
				rf.mu.Unlock()
				return
			}
			args := AppendEntriesArgs{
				Term:         rf.Term(),
				LeaderId:     rf.me,
				PrevLogIndex: rf.state.nextIndex[peer] - 1,
				PrevLogTerm:  rf.safeTerm(rf.state.nextIndex[peer] - 1),
				Entries:      rf.safeEntries(rf.state.nextIndex[peer]),
				LeaderCommit: rf.state.commitIndex,
			}
			rf.mu.Unlock()

			reply := AppendEntriesReply{}
			if rf.appendEntries(peer, &args, &reply) {
				rf.mu.Lock()
				defer rf.mu.Unlock()
				defer T(time.Now(), "%v: AppendEntriesReply(from: %d, to: %d)", &rf.state, peer, rf.me)

				if reply.Term > rf.Term() {
					rf.state.ToFollower(reply.Term)
				}

				if reply.Success {
					rf.state.nextIndex[peer] = rf.state.LastLogIndex() + 1
					rf.state.matchIndex[peer] = rf.state.LastLogIndex()
				} else {
					rf.state.nextIndex[peer]--
				}

				rf.updateCommit()
			}
		}(i)
	}
}

func (rf *Raft) updateCommit() {
	for N := rf.state.LastLogIndex(); N > rf.state.commitIndex; N-- {
		majority := 1
		for i := range rf.peers {
			if i == rf.me {
				continue
			}

			if rf.state.matchIndex[i] >= N {
				majority++
			}
		}

		if majority > len(rf.peers)/2 {
			rf.state.commitIndex = N
			go rf.Apply()
			return
		}
	}
}

func (rf *Raft) RunLeader() {
	for !rf.killed() {
		rf.mu.Lock()
		if rf.state.state != LeaderState {
			rf.mu.Unlock()
			break
		}
		rf.mu.Unlock()
		rf.BroadcastHeartbeat()
		time.Sleep(50 * time.Millisecond)
	}
}
