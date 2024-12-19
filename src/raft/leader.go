package raft

import "time"

func (rf *Raft) safeTerm(index int) int {
	if index < 0 || index > len(rf.state.log)-1 {
		return -1
	}
	return rf.state.log[index].Term
}

func (rf *Raft) safeEntries(index int) []LogEntry {
	if index < 0 || index > len(rf.state.log)-1 {
		return make([]LogEntry, 0)
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
				Entries:      rf.safeEntries(rf.state.nextIndex[peer] - 1),
				LeaderCommit: rf.state.commitIndex,
			}
			rf.mu.Unlock()

			reply := AppendEntriesReply{}
			if rf.appendEntries(peer, &args, &reply) {
				rf.mu.Lock()
				defer rf.mu.Unlock()

				if reply.Term > rf.Term() {
					rf.state.ToFollower(reply.Term)
				}

				if reply.Success {
					rf.state.nextIndex[peer] = args.PrevLogIndex + len(args.Entries) + 1
					rf.state.matchIndex[peer] = args.PrevLogIndex + len(args.Entries)
				} else {
					rf.state.nextIndex[peer]--
				}
			}
		}(i)
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
