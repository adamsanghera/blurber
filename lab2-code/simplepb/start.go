package simplepb

import (
	"log"
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
func (srv *PBServer) Start(command interface{}) (
	index int, view int, ok bool) {
	srv.mu.Lock()
	// log.Printf("Server %d: Received start RPC\n", srv.me)
	// do not process command if status is not NORMAL
	// and if i am not the primary in the current view
	if srv.status != NORMAL {
		return -1, srv.currentView, false
	} else if GetPrimary(srv.currentView, len(srv.peers)) != srv.me {
		return -1, srv.currentView, false
	}

	// At this point, we are a primary in a normal state

	// Append the command, update state
	srv.log = append(srv.log, command)
	index = len(srv.log) - 1
	view = srv.currentView
	commit := srv.commitIndex

	log.Printf("PRIMARY: Start RPC releasing lock, with state values {index: %d} {view: %d} {commit: %d}\n", index, view, commit)

	srv.mu.Unlock()

	// Propagate this action to our backups.
	// This needs to happen asynchronously, so we need a sub-routine here.
	// TODO Write subroutine for propagation:
	//   One caller worker for each backup, takes commands as inbox
	//   Each caller sends bool to outbox, which is processed by a single routine.
	//   If the receiver observes that all commits succeeded, then it takes the mutex
	//     and updates the commit-index.
	// go func() {
	// 	outbox := make(chan bool)

	// 	go func() {
	// 		log.Printf("PRIMARY: ACC for %d: Launch\n", index)
	// 		// Start with 1, because the master has prepared it.
	// 		positivePreps := 1
	// 		negativePreps := 0

	// 		// Continue to receive messages until we get all callbacks
	// 		for (positivePreps + negativePreps) < len(srv.peers) {
	// 			isGood := <-outbox
	// 			log.Printf("PRIMARY: ACC for %d: received response num %d: %v\n", index, positivePreps+negativePreps, isGood)
	// 			if isGood {
	// 				positivePreps++
	// 			} else {
	// 				negativePreps++
	// 			}

	// 			// If positive prepare messages are in majority.
	// 			if positivePreps > len(srv.peers)/2 {
	// 				// Committed
	// 				srv.mu.Lock()
	// 				if srv.commitIndex < index {
	// 					log.Printf("PRIMARY: ACC for %d: Updating commit idx to %d\n", index, index)
	// 					srv.commitIndex = index
	// 				}
	// 				srv.mu.Unlock()
	// 			}
	// 		}

	// 		log.Printf("PRIMARY: ACC for %d: All messages received. Closing outbox...\n", index)

	// 		// At this point, all preparers have sent their message
	// 		close(outbox)
	// 	}()

	// 	for i := range srv.peers {
	// 		// log.Printf("PRIMARY: Woohoo view %d, for PRIMARY of %d\n", view, i, len(srv.peers))
	// 		// log.Printf("PRIMARY: Result of GetPrimary: %d", GetPrimary(srv.currentView, len(srv.peers)))
	// 		if GetPrimary(srv.currentView, len(srv.peers)) != i {
	// 			go func(backup *labrpc.ClientEnd, backupNum int) {
	// 				pArgs := &PrepareArgs{
	// 					View:          view,
	// 					PrimaryCommit: commit,
	// 					Index:         index,
	// 					Entry:         command,
	// 				}
	// 				pReply := &PrepareReply{}
	// 				log.Printf("PRIMARY: PREP for %d -> %d: Sending RPC\n", index, backupNum)
	// 				ok := backup.Call("PBServer.Prepare", pArgs, pReply)
	// 				log.Printf("PRIMARY: PREP for %d -> %d: Response recvd: {%v}\n", index, backupNum, ok)
	// 				if !ok {
	// 					outbox <- false
	// 				} else {
	// 					// ADAM: Do we need to use the View element of the reply at all?
	// 					//   --> We probably need it to send commits.
	// 					log.Printf("PRIMARY: PREP for %d -> %d: Forwarding to ACC\n", index, backupNum)
	// 					outbox <- pReply.Success
	// 				}
	// 			}(srv.peers[i], i)
	// 		}
	// 	}
	// }()

	go srv.syncrhonize(index, view, commit, command)

	// Your code here

	return index, view, true
}
