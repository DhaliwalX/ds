package raft

func (rf *Raft) Apply() {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if rf.state.commitIndex <= rf.state.lastApplied {
		return
	}

	for i := rf.state.lastApplied + 1; i <= rf.state.commitIndex; i++ {
		rf.applyCh <- ApplyMsg{
			CommandValid: true,
			Command:      rf.state.log[i].Command,
			CommandIndex: i,
		}
	}

	rf.state.lastApplied = rf.state.commitIndex
}
