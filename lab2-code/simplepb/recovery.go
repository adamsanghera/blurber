package simplepb

import (
	"context"
	"log"

	"github.com/adamsanghera/blurber-protobufs/dist/replication"
)

// Recovery is the RPC handler for the Recovery RPC
func (srv *PBServer) Recovery(args *replication.RecoveryArgs, reply *replication.RecoveryReply) {
	// Your code here
	srv.mu.Lock()
	log.Printf("PRIMARY: RECOVERY for %d: Launched\n", args.Server)

	if args.View != srv.currentView {
		log.Printf("PRIMARY: RECOVERY for %d: View out of sync; this %d vs req %d\n", args.Server, srv.currentView, args.View)
	} else {
		log.Printf("PRIMARY: RECOVERY for %d: Acking positively\n", args.Server)
		reply.Entries = srv.log
		reply.PrimaryCommit = srv.commitIndex
		reply.View = srv.currentView
		reply.Success = true
	}
	srv.mu.Unlock()
	log.Printf("PRIMARY: RECOVERY for %d: Returning\n", args.Server)
	return
}

func (srv *PBServer) sendRecovery() {
	good := false

	srv.mu.Lock()
	log.Printf("Server %d: CPP in REC: Beginning recovery, locked until completion...\n", srv.me)
	srv.status = RECOVERING
	arg := &replication.RecoveryArgs{
		View:   srv.currentView,
		Server: srv.me,
	}

	for !good {
		log.Printf("Server %d: CPP in REC: Sending recovery request\n", srv.me)
		rep, err := srv.peers[GetPrimary(srv.currentView, int32(len(srv.peers)))].Recovery(context.Background(), arg)
		good = err == nil

		if good {
			log.Printf("Server %d: CPP in REC: Got a response from primary!\n", srv.me)
			good = rep.Success
			if good {
				log.Printf("Server %d: CPP in REC: Accepted for recovery!\n", srv.me)
				srv.log = rep.Entries
				srv.commitIndex = rep.PrimaryCommit
				srv.currentView = rep.View
				srv.status = NORMAL
			}
		}
	}

	srv.mu.Unlock()
}
