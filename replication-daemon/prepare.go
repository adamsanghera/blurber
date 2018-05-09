package pbdaemon

import (
	"context"
	"log"

	"github.com/adamsanghera/blurber-protobufs/dist/replication"
)

// Prepare is the RPC handler for the Prepare RPC
func (srv *PBServer) Prepare(ctx context.Context, args *replication.PrepareArgs) (*replication.PrepareReply, error) {
	log.Printf("Server %d: Prep starting on {%d}\n", srv.me, args.Index)

	if srv.status != NORMAL {
		return &replication.PrepareReply{
			Success: false,
			View:    srv.currentView,
		}, nil
	}

	callback := make(chan bool)

	log.Printf("Server %d: Prep sending {%d} to CPP\n", srv.me, args.Index)
	srv.prepChan <- &callbackArg{
		callback: callback,
		args:     args,
		handled:  false,
	}

	reply := &replication.PrepareReply{
		Success: <-callback,
		View:    srv.currentView,
	}

	log.Printf("Server %d: Prep ack {%v} from CPP for message of idx {%d}, returning\n", srv.me, reply.Success, args.Index)

	close(callback)

	return reply, nil
}
