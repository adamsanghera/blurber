package simplepb

import (
	"log"

	"github.com/adamsanghera/blurber-protobufs/dist/replication"
)

// Prepare is the RPC handler for the Prepare RPC
func (srv *PBServer) Prepare(args *replication.PrepareArgs, reply *replication.PrepareReply) {
	log.Printf("Server %d: Prep starting on {%d}\n", srv.me, args.Index)
	// Your code here
	// we need to pass in this prep to a centralized prepare-processor routine.
	// That routine will take args, and return bools.
	// The Prepare function itself blocks until that routine sends a message back.
	//
	// This behavior implies that we have to have some global-level synchronization, which the preparer routine can
	//   use to communicate with the preparing routines.
	//   Perhaps, on receipt, this function will create an entry in a global hash-map, of int -> channel.
	//     Then, this function will push its args into a channel.
	//   The routine takes these args, and processes them in its own way.
	//   When an arg is processed, then it sends a process-success message to the channel in the global hash map, and closes the channel.
	//     This function blocks on receiving from that channel.
	//     It cleans up the hash map entry, and then returns the prep-success message.

	if srv.status != NORMAL {
		reply.Success = false
		reply.View = srv.currentView
		return
	}

	callback := make(chan bool)

	log.Printf("Server %d: Prep sending {%d} to CPP\n", srv.me, args.Index)
	srv.prepChan <- &callbackArg{
		callback: callback,
		args:     args,
		handled:  false,
	}

	reply.Success = <-callback
	reply.View = srv.currentView

	log.Printf("Server %d: Prep ack {%v} from CPP for message of idx {%d}, returning\n", srv.me, reply.Success, args.Index)

	close(callback)

	return
}
