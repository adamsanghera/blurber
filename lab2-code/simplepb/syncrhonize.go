package simplepb

import (
	"log"
)

func (srv *PBServer) syncrhonize(index int, view int, commit int, command interface{}) {
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
		if GetPrimary(srv.currentView, len(srv.peers)) != i {
			go func(backupAddr string, backupNum int) {
				pArgs := &PrepareArgs{
					View:          view,
					PrimaryCommit: commit,
					Index:         index,
					Entry:         command,
				}
				pReply := &PrepareReply{}
				log.Printf("PRIMARY: PREP for %d -> %d: Sending RPC\n", index, backupNum)
				ok := backupAddr.Call("PBServer.Prepare", pArgs, pReply)
				log.Printf("PRIMARY: PREP for %d -> %d: Response recvd: {%v}\n", index, backupNum, ok)
				if !ok {
					outbox <- false
				} else {
					// ADAM: Do we need to use the View element of the reply at all?
					//   --> We probably need it to send commits.
					log.Printf("PRIMARY: PREP for %d -> %d: Forwarding to ACC\n", index, backupNum)
					outbox <- pReply.Success
				}
			}(srv.peers[i], i)
		}
	}
}
