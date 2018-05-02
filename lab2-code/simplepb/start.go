package simplepb

import (
	"log"

	"github.com/adamsanghera/blurber-protobufs/dist/replication"
)

// Start() is invoked by tester on some replica server to replicate a
// command.  Only the primary should process this request by appending
// the command to its log and then return *immediately* (while the log is being replicated to backup servers).
// if this server isn't the primary, returns false.
// Note that since the function returns immediately, there is no guarantee that this command
// will ever be committed upon return, since the primary
// may subsequently fail before replicating the command to all servers
//
// The first return value is the index that the command will appear at
// *if it's eventually committed*. The second return value is the current
// view. The third return value is true if this server believes it is
// the primary.
func (srv *PBServer) Start(command *replication.Command) (
	index int32, view int32, ok bool) {
	srv.mu.Lock()
	// log.Printf("Server %d: Received start RPC\n", srv.me)
	// do not process command if status is not NORMAL
	// and if i am not the primary in the current view
	if srv.status != NORMAL {
		return -1, srv.currentView, false
	} else if GetPrimary(srv.currentView, int32(len(srv.peers))) != srv.me {
		return -1, srv.currentView, false
	}

	// At this point, we are a primary in a normal state

	// Append the command, update state
	srv.log = append(srv.log, command)
	index = int32(len(srv.log) - 1)
	view = srv.currentView
	commit := srv.commitIndex

	log.Printf("PRIMARY: Start RPC releasing lock, with state values {index: %d} {view: %d} {commit: %d}\n", index, view, commit)

	srv.mu.Unlock()

	go srv.syncrhonize(index, view, commit, command)

	// Your code here

	return index, view, true
}
