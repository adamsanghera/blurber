package pbdaemon

import (
	"context"
	"log"

	"github.com/adamsanghera/blurber-protobufs/dist/common"
	"github.com/adamsanghera/blurber-protobufs/dist/replication"
	"google.golang.org/grpc"
)

func (srv *PBServer) connectPeerUnsafe(addr string) error {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Printf("ReplicationD: Failed to connect to the Replication Daemon at %s", addr)
		return err
	}
	newClient := replication.NewReplicationClient(conn)
	srv.peers = append(srv.peers, newClient)
	srv.peerAddresses = append(srv.peerAddresses, addr)

	return nil
}

// SharePeers is an RPC function for spreading machine addresses within the network.
// It is only ever invoked by the Leader.  Top-down coercion at its finest.
// Every node in the network must maintain consistent lists of peers, otherwise
// they will eventually disagree on who the leader of a given view is – which is bad!
func (srv *PBServer) SharePeers(ctx context.Context, in *replication.Peers) (*common.Empty, error) {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	// Remember where we came from
	myAddr := srv.peerAddresses[srv.me]

	srv.peerAddresses = make([]string, 0)
	srv.peers = make([]replication.ReplicationClient, 0)

	for idx, addr := range in.Addresses {
		if addr == myAddr {
			srv.me = int32(idx)
		}

		srv.connectPeerUnsafe(addr)
	}

	return &common.Empty{}, nil
}

// GetLeaderInfo is a safe function for getting a snapshot of who the replication daemon believes is its leader
func (srv *PBServer) GetLeaderInfo() (int32, string) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	idx := GetPrimary(srv.currentView, int32(len(srv.peers)))
	return idx, srv.peerAddresses[idx]
}

func (srv *PBServer) GetView() int32 {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	return srv.currentView
}
