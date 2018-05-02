package simplepb

import (
	"context"
	"log"

	"github.com/adamsanghera/blurber-protobufs/dist/replication"
)

func (srv *PBServer) syncrhonize(index int32, view int32, commit int32, command *replication.Command) {
	outbox := make(chan bool)

	go func() {
		log.Printf("PRIMARY: ACC for %d: Launch\n", index)
		// Start with 1, because the master has prepared it.
		positivePreps := 1
		negativePreps := 0

		// Continue to receive messages until we get all callbacks
		for (positivePreps + negativePreps) < len(srv.peers) {
			isGood := <-outbox
			log.Printf("PRIMARY: ACC for %d: received response num %d: %v\n", index, positivePreps+negativePreps, isGood)
			if isGood {
				positivePreps++
			} else {
				negativePreps++
			}

			// If positive prepare messages are in majority.
			if positivePreps > len(srv.peers)/2 {
				// Committed
				srv.mu.Lock()
				if srv.commitIndex < index {
					log.Printf("PRIMARY: ACC for %d: Updating commit idx to %d\n", index, index)
					srv.commitIndex = index
				}
				srv.mu.Unlock()
			}
		}

		log.Printf("PRIMARY: ACC for %d: All messages received. Closing outbox...\n", index)

		// At this point, all preparers have sent their message
		close(outbox)
		if !(positivePreps > negativePreps) {
			srv.syncrhonize(index, view, commit, command)
		}
	}()

	for i := range srv.peers {
		// log.Printf("PRIMARY: Woohoo view %d, for PRIMARY of %d\n", view, i, len(srv.peers))
		// log.Printf("PRIMARY: Result of GetPrimary: %d", GetPrimary(srv.currentView, len(srv.peers)))
		if GetPrimary(srv.currentView, int32(len(srv.peers))) != int32(i) {
			go func(backup replication.ReplicationClient, backupNum int) {
				pArgs := &replication.PrepareArgs{
					View:          view,
					PrimaryCommit: commit,
					Index:         index,
					Entry:         command,
				}
				pReply := &replication.PrepareReply{}
				log.Printf("PRIMARY: PREP for %d -> %d: Sending RPC\n", index, backupNum)
				pReply, err := backup.Prepare(context.Background(), pArgs)
				log.Printf("PRIMARY: PREP for %d -> %d: Response recvd: {%v}\n", index, backupNum, err == nil)
				if err != nil {
					outbox <- false
				} else {
					log.Printf("PRIMARY: PREP for %d -> %d: Forwarding to ACC\n", index, backupNum)
					outbox <- pReply.Success
				}
			}(srv.peers[i], i)
		}
	}
}
