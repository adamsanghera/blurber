package simplepb

import (
	"log"

	"golang.org/x/net/context"

	"github.com/adamsanghera/blurber-protobufs/dist/replication"
)

// Recovery is the RPC handler for the Recovery RPC
func (srv *PBServer) Recovery(ctx context.Context, args *replication.RecoveryArgs) (*replication.RecoveryReply, error) {
	// Your code here
	srv.mu.Lock()
	defer srv.mu.Unlock()

	peerIdx := -1
	for i, a := range srv.peerAddresses {
		if a == args.Address {
			peerIdx = i
		}
	}

	// Received contact from a new machine, need to store new connection
	if peerIdx == -1 {
		log.Printf("PRIMARY: Recovery for %s: Bootstrapping", args.Address)
		err := srv.connectPeer(args.Address)
		if err != nil {
			log.Printf("Failed to make client connection with peer")
		}
		peerIdx = len(srv.peerAddresses) - 1
	}

	log.Printf("PRIMARY: Recovery for %s: Sending reply", args.Address)

	if args.View != srv.currentView {
		log.Printf("PRIMARY: RECOVERY for %s: View out of sync; this %d vs req %d\n", args.Address, srv.currentView, args.View)
	} else {
		log.Printf("PRIMARY: RECOVERY for %s: Acking positively\n", args.Address)
		return &replication.RecoveryReply{
			View:          srv.currentView,
			Entries:       srv.log,
			PrimaryCommit: srv.commitIndex,
			Success:       true,
			Peers:         srv.peerAddresses,
			NewIdentity:   int32(peerIdx),
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
		View:    srv.currentView,
		Address: srv.peerAddresses[srv.me],
	}

	for !good {
		log.Printf("Server %d: CPP in REC: Sending recovery request from addr %s\n", srv.me, srv.peerAddresses[srv.me])

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
				srv.me = rep.NewIdentity

				// Re-forge connections to other Replication Daemons.
				srv.peerAddresses = make([]string, 0)
				srv.peers = make([]replication.ReplicationClient, 0)
				for _, addr := range rep.Peers {
					srv.connectPeer(addr)
				}

				for idx := int32(1); idx <= srv.commitIndex; idx++ {
					srv.commitChan <- srv.log[idx]
				}
			}
		}
	}

	srv.mu.Unlock()
}
