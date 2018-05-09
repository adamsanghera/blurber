package pbdaemon

import (
	"log"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/any"

	"github.com/adamsanghera/blurber-protobufs/dist/replication"
)

func TestSharePeers(t *testing.T) {
	type args struct {
		thisAddress   string
		leaderAddress string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		struct {
			name string
			args args
		}{
			"hi",
			args{
				thisAddress:   "127.0.0.1:4001",
				leaderAddress: "127.0.0.1:4000",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewReplicationDaemon(tt.args.leaderAddress, tt.args.leaderAddress)
			if srv.me != 0 {
				t.Fatalf("Incorrect servery entry")
			}

			if srv.currentView != 0 {
				t.Fatalf("Incorrect view")
			}

			if len(srv.peers) != 1 {
				t.Fatalf("Incorrect peers list size")
			}

			srv.Replicate(&replication.Command{
				Cmd: &any.Any{
					TypeUrl: "Add",
					Value:   []byte("adam"),
				},
			})

			time.Sleep(10 * time.Millisecond)

			if len(srv.log) != 2 {
				t.Fatalf("Failed to store commands")
			}

			follower := NewReplicationDaemon(tt.args.thisAddress, tt.args.leaderAddress)
			time.Sleep(10 * time.Millisecond)
			follower2 := NewReplicationDaemon("127.0.0.1:4002", tt.args.leaderAddress)

			time.Sleep(10 * time.Millisecond)

			srv.mu.Lock()
			follower.mu.Lock()
			follower2.mu.Lock()

			log.Printf("Leader peers: %v", srv.peerAddresses)
			log.Printf("Follower peers: %v", follower.peerAddresses)
			log.Printf("Follower2 peers: %v", follower2.peerAddresses)

			for idx, addr := range srv.peerAddresses {
				if addr != follower.peerAddresses[idx] {
					t.Fatalf("Follower does not have up to date peers list %v \n vs %v", srv.peerAddresses, follower.peerAddresses)
				}
				if addr != follower2.peerAddresses[idx] {
					t.Fatalf("Follower does not have up to date peers list %v \n vs %v", srv.peerAddresses, follower2.peerAddresses)
				}
			}
			srv.mu.Unlock()
			follower.mu.Unlock()
			follower2.mu.Unlock()
		})
	}
}
