package pbdaemon

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

	if GetPrimary(srv.currentView, int32(len(srv.peerAddresses))) != srv.me {
		return &replication.RecoveryReply{
			View:    srv.currentView,
			Peers:   srv.peerAddresses,
			Success: false,
		}, nil
	}

	peerIdx := -1
	for i, a := range srv.peerAddresses {
		if a == args.Address {
			peerIdx = i
		}
	}

	// Received contact from a new machine, need to store new connection
	if peerIdx == -1 {
		log.Printf("PRIMARY: Recovery for %s: Bootstrapping", args.Address)
		err := srv.connectPeerUnsafe(args.Address)
		if err != nil {
			log.Printf("Failed to make client connection with peer")
		}
		peerIdx = len(srv.peerAddresses) - 1
	}

	for idx, client := range srv.peers[:len(srv.peers)-1] {
		if idx != int(srv.me) {
			log.Printf("PRIMARY: Sending peer list to %s", srv.peerAddresses[idx])

			// Spawn a process to share peers with all our buddies
			go func(client replication.ReplicationClient) {
				_, err := client.SharePeers(context.Background(), &replication.Peers{Addresses: srv.peerAddresses})
				if err != nil {
					log.Printf("PRIMARY: Failed to connect to client %s", srv.peerAddresses[idx])
				}
			}(client)
		}
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

// sendRecovery is a helper function for initiating a daemon's recovery.
// It triggers the Recovery RPC, targeting the primary.
func (srv *PBServer) sendRecovery() {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	good := false

	log.Printf("Server %d: CPP in REC: Beginning recovery, locked until completion...\n", srv.me)
	srv.status = RECOVERING
	arg := &replication.RecoveryArgs{
		View:    srv.currentView,
		Address: srv.peerAddresses[srv.me],
	}

	for !good {
		log.Printf("Server %d: CPP in REC: Sending recovery request from addr %s\n", srv.me, srv.peerAddresses[srv.me])

		leaderIdx := GetPrimary(srv.currentView, int32(len(srv.peers)))

		log.Printf("Sending recovery to %s", srv.peerAddresses[leaderIdx])

		rep, err := srv.peers[leaderIdx].Recovery(context.Background(), arg)
		good = (err == nil)

		// Made contact with sooooomebody
		if good {
			log.Printf("Server %d: CPP in REC: Got a response from primary!\n", srv.me)
			good = rep.Success

			// Made contact with the one true leader
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
					srv.connectPeerUnsafe(addr)
				}

				// Throw all of the new commands down the hatch
				for idx := int32(1); idx <= srv.commitIndex; idx++ {
					srv.CommitChan <- srv.log[idx]
				}
			} else {
				log.Printf("Server %d: CPP in REC: Got a response from non-primary!  Redirecting to primary...\n", srv.me)
				// This should redirect us to the One True Leader
				srv.peerAddresses = rep.Peers
				srv.currentView = rep.View
			}
		}
	}
}
