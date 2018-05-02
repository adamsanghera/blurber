package simplepb

import (
	"log"

	"github.com/adamsanghera/blurber-protobufs/dist/replication"
)

// Start is the entry point to our replication engine
func (srv *PBServer) Start(command *replication.Command) (
	index int32, view int32, ok bool) {
	srv.mu.Lock()

	// If we're panicking or not the primary, quit quit quit
	if srv.status != NORMAL {
		srv.mu.Unlock()
		return -1, srv.currentView, false
	} else if GetPrimary(srv.currentView, int32(len(srv.peers))) != srv.me {
		srv.mu.Unlock()
		return -1, srv.currentView, false
	}

	// Only primaries shall pass...
	srv.log = append(srv.log, command)
	index = int32(len(srv.log) - 1)
	view = srv.currentView
	commit := srv.commitIndex

	log.Printf("PRIMARY: Start RPC releasing lock, with state values {index: %d} {view: %d} {commit: %d}\n", index, view, commit)
	srv.mu.Unlock()

	go srv.syncrhonize(index, view, commit, command)

	return index, view, true
}
