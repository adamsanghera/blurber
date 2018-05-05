package simplepb

import (
	"context"
	"log"

	"github.com/adamsanghera/blurber-protobufs/dist/replication"
)

// Recovery is the RPC handler for the Recovery RPC
func (srv *PBServer) Recovery(ctx context.Context, args *replication.RecoveryArgs) (*replication.RecoveryReply, error) {
	// Your code here
	srv.mu.Lock()
	defer srv.mu.Unlock()

	addr, ok := ctx.Value("address").(string)
	if !ok {
		// do smth
		log.Printf("RECOVERY: Insufficient context in recovery request... no address provided")
	}

	known := false
	for _, a := range srv.peerAddresses {
		known = known || a == addr
	}

	// New machine, needs to store new connection
	if !known {
		err := srv.connectPeer(addr)
		if err != nil {
			log.Printf("Failed to make client connection with peer")
		}
	}

	log.Printf("PRIMARY: RECOVERY for %d: Launched\n", args.Server)

	if args.View != srv.currentView {
		log.Printf("PRIMARY: RECOVERY for %d: View out of sync; this %d vs req %d\n", args.Server, srv.currentView, args.View)
	} else {
		log.Printf("PRIMARY: RECOVERY for %d: Acking positively\n", args.Server)
		return &replication.RecoveryReply{
			Entries:       srv.log,
			PrimaryCommit: srv.commitIndex,
			View:          srv.currentView,
			Success:       true,
		}, nil
	}
	return &replication.RecoveryReply{}, nil
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
		ctx := context.WithValue(context.Background(), "address", srv.peerAddresses[srv.me])

		rep, err := srv.peers[GetPrimary(srv.currentView, int32(len(srv.peers)))].Recovery(ctx, arg)
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
