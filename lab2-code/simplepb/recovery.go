package simplepb

import (
	"context"
	"log"

	"github.com/adamsanghera/blurber-protobufs/dist/replication"
	"google.golang.org/grpc"
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

	idx, ok := ctx.Value("index").(int32)
	if !ok {
		// do smth
		log.Printf("RECOVERY: Insufficient context in recovery request... no index provided")
	}

	// New machine, needs to be bootstrapped
	if idx == -1 {
		connection, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			log.Printf("RECOVERY: Unable to establish a new client connection")
		}
		srv.peers = append(srv.peers, replication.NewReplicationClient(connection))
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
