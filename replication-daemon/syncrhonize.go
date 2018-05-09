package pbdaemon

import (
	"context"
	"log"

	"github.com/adamsanghera/blurber-protobufs/dist/replication"
)

// synchronize is a background-routine, used to synchronize follower logs with the leader log
// synchronize is a composite of two functions.
// 1. A receiver function, that listens for responses to prepare RPCs
// 2. A sender function, which propagates the new to peers
func (srv *PBServer) syncrhonize(index int32, view int32, commit int32, command *replication.Command) {
	outbox := make(chan *replication.PrepareReply)

	go func() {
		log.Printf("PRIMARY: ACC for %d: Launch\n", index)

		// positivePreps starts at 1, because the master already accepts it
		positivePreps := 1
		negativePreps := 0

		// Continue to wait on messages until all callbacks reply
		for (positivePreps + negativePreps) < len(srv.peers) {
			reply := <-outbox

			log.Printf("PRIMARY: ACC for %d: received response num %d: %v\n", index, positivePreps+negativePreps, reply.Success)

			if reply.Success {
				positivePreps++
			} else {
				// Out of sync with network, need to recover
				if srv.currentView < reply.View {
					defer srv.sendRecovery()
				} else {
					// failed to connect, or anomalous reject, but still in sync – keep going
					negativePreps++
				}
			}
		}

		log.Printf("PRIMARY: ACC for %d: Done receiving messages. Closing outbox...\n", index)
		close(outbox)

		// Positive prepares are majority.
		srv.mu.Lock()
		if positivePreps > len(srv.peers)/2 {
			// Committed
			if srv.commitIndex < index {
				log.Printf("PRIMARY: ACC for %d: Updating commit idx to %d\n", index, index)
				for idx := srv.commitIndex + 1; idx <= index; idx++ {
					srv.CommitChan <- srv.log[idx]
				}
				srv.commitIndex = index
				srv.PropagateCommitsUnsafe()
			}
		}
		srv.mu.Unlock()

		// If we failed to reach consensus, but have not fallen behind the view, we retry synchronizing.
		if !(positivePreps > negativePreps) {
			srv.syncrhonize(index, view, commit, command)
		}
	}()

	// Spawn a routine for every peer, not including self
	srv.mu.Lock()
	for i := range srv.peers {
		if GetPrimary(srv.currentView, int32(len(srv.peers))) != int32(i) {
			go func(backup replication.ReplicationClient, backupNum int) {
				pArgs := &replication.PrepareArgs{
					View:          view,
					PrimaryCommit: commit,
					Index:         index,
					Entry:         command,
				}

				log.Printf("PRIMARY: PREP for %d -> %d: Sending RPC\n", index, backupNum)
				pReply, err := backup.Prepare(context.Background(), pArgs)

				if err != nil {
					log.Printf("PRIMARY: PREP for %d -> %d: No response\n", index, backupNum)
					outbox <- &replication.PrepareReply{Success: false}
				} else {
					log.Printf("PRIMARY: PREP for %d -> %d: Good response\n", index, backupNum)
					outbox <- pReply
				}
			}(srv.peers[i], i)
		}
	}
	srv.mu.Unlock()
}
