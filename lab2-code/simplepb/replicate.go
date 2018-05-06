package simplepb

import (
	"context"
	"log"

	"github.com/adamsanghera/blurber-protobufs/dist/replication"
)

// Replicate is the entry point to our replication engine
func (srv *PBServer) Replicate(ctx context.Context, command *replication.Command) (*replication.ReplicateReply, error) {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	// If we're panicking or not the primary, quit quit quit
	if srv.status != NORMAL {
		return &replication.ReplicateReply{
			Index:   -1,
			View:    srv.currentView,
			Success: false,
		}, nil
	} else if GetPrimary(srv.currentView, int32(len(srv.peers))) != srv.me {
		return &replication.ReplicateReply{
			Index:   -1,
			View:    srv.currentView,
			Success: false,
		}, nil
	}

	// Only primaries shall pass...
	srv.log = append(srv.log, command)
	reply := &replication.ReplicateReply{
		Index:   int32(len(srv.log) - 1),
		View:    srv.currentView,
		Success: true,
	}

	commit := srv.commitIndex

	log.Printf("PRIMARY: Replicate releasing lock, with state values {index: %d} {view: %d} {commit: %d}\n", reply.Index, reply.View, commit)

	go srv.syncrhonize(reply.Index, reply.View, commit, command)

	return reply, nil
}
