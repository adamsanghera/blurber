package pbdaemon

import (
	"errors"
	"log"

	"github.com/adamsanghera/blurber-protobufs/dist/replication"
)

// Replicate is the entry point to our replication engine
func (srv *PBServer) Replicate(cmd *replication.Command) error {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	// If we're panicking or not the primary, quit quit quit
	if srv.status != NORMAL {
		return errors.New("Replication Daemon is not OK right now")
	} else if GetPrimary(srv.currentView, int32(len(srv.peers))) != srv.me {
		return errors.New("Replication Daemon is *not* the master daemon")
	}

	// Only primaries shall pass...
	srv.log = append(srv.log, cmd)
	commit := srv.commitIndex

	log.Printf("PRIMARY: Replicate releasing lock, with state values {index: %d} {view: %d} {commit: %d}\n", reply.Index, reply.View, commit)

	go srv.syncrhonize(int32(len(srv.log)-1), srv.currentView, commit, cmd)

	return nil
}
