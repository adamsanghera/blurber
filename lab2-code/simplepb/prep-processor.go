package simplepb

import (
	"log"
	"sort"
)

// prepareProcessor is the heart of the replication daemon
// The processor waits for commands, and is guaranteed to process them sequentially.
// The processor also detects when it falls behind, and initiates recovery in that case.
func (srv *PBServer) prepareProcessor() {
	recvdArgs := make([]*callbackArg, 0)
	log.Printf("Server %d: Beginning Central Prep Processor\n", srv.me)

	for callArg := range srv.prepChan {
		log.Printf("Server %d: CPP: Received arg {%d}\n", srv.me, callArg.args.Index)

		// Check if recovery is necessary
		if callArg.args.PrimaryCommit >= int32(len(srv.log)) || srv.currentView < callArg.args.View {

			log.Printf("Server %d: CPP: This log {%d} is behind primary {%d}!\n", srv.me, srv.commitIndex, callArg.args.PrimaryCommit)

			callArg.callback <- false
			for _, arg := range recvdArgs {
				arg.callback <- false
			}
			recvdArgs = make([]*callbackArg, 0)

			srv.sendRecovery()

		} else {
			// Check if it's just a commitIdx update
			log.Printf("Server %d: CPP: appending {%d} to backlog\n", srv.me, callArg.args.Index)
			// Append the received prepare to the log
			recvdArgs = append(recvdArgs, callArg)

			// Sort the log
			sort.Slice(recvdArgs, func(i, j int) bool {
				return recvdArgs[i].args.Index < recvdArgs[j].args.Index
			})

			log.Printf("Server %d: CPP: after sorting, backlog is %v\n", srv.me, recvdArgs[0].args.Index)
			// While the youngest argument is the next index we're looking for
			srv.mu.Lock()
			for len(recvdArgs) > 0 && recvdArgs[0].args.Index == int32(len(srv.log)) {
				// Append cmd to log, and send an okay to the prep RPC
				log.Printf("Server %d: CPP: processing arg of idx %d\n", srv.me, recvdArgs[0].args.Index)
				srv.log = append(srv.log, recvdArgs[0].args.Entry)
				if srv.commitIndex != recvdArgs[0].args.PrimaryCommit {
					for idx := srv.commitIndex + 1; idx <= recvdArgs[0].args.PrimaryCommit; idx++ {
						srv.commitChan <- srv.log[idx]
					}
				}
				srv.commitIndex = recvdArgs[0].args.PrimaryCommit
				log.Printf("Server %d: CPP: sending good news to callback...\n", srv.me)
				recvdArgs[0].callback <- true
				recvdArgs = recvdArgs[1:]
			}
			srv.mu.Unlock()

		}
	}

	log.Printf("Server %d: CPP: received close... shutting down\n", srv.me)
}
