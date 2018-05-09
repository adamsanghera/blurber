package pbdaemon

import (
	"log"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/any"

	"github.com/adamsanghera/blurber-protobufs/dist/replication"
)

func TestNewReplicationDaemon(t *testing.T) {
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
			srv.Replicate(&replication.Command{
				Cmd: &any.Any{
					TypeUrl: "Add",
					Value:   []byte("bob"),
				},
			})
			srv.Replicate(&replication.Command{
				Cmd: &any.Any{
					TypeUrl: "Add",
					Value:   []byte("bob"),
				},
			})

			time.Sleep(1 * time.Second)

			if len(srv.log) != 4 {
				t.Fatalf("Failed to store commands")
			}

			follower := NewReplicationDaemon(tt.args.thisAddress, tt.args.leaderAddress)
			time.Sleep(1 * time.Second)
			follower2 := NewReplicationDaemon("127.0.0.1:4002", tt.args.leaderAddress)

			time.Sleep(2 * time.Second)

			if len(follower.log) != len(srv.log) {
				t.Fatalf("Follower failed to receive replicated log\n%v\nvs\n%v", follower.log, srv.log)
			}

			if len(follower2.peerAddresses) != len(srv.peerAddresses) {
				t.Fatalf("Failed to replicate peer list:\n\t{%v}\n\tvs\n\t{%v}", follower2.peerAddresses, srv.peerAddresses)
			}

			srv.Replicate(&replication.Command{
				Cmd: &any.Any{
					TypeUrl: "Sub",
					Value:   []byte("adam"),
				},
			})
			srv.Replicate(&replication.Command{
				Cmd: &any.Any{
					TypeUrl: "Sub",
					Value:   []byte("bob"),
				},
			})
			srv.Replicate(&replication.Command{
				Cmd: &any.Any{
					TypeUrl: "Sub",
					Value:   []byte("bob"),
				},
			})

			time.Sleep(1 * time.Second)

			if len(follower.log) != len(srv.log) {
				log.Printf("{%v}", follower.log)
				log.Printf("{%v}", srv.log)
				t.Fatalf("Failed to replicate log to follower")
			}

			if srv.peerAddresses[srv.me] != tt.args.leaderAddress {
				t.Fatalf("Mistaken identity %d of %v", srv.me, srv.peerAddresses)
			}

			if follower.peerAddresses[follower.me] != tt.args.thisAddress {
				t.Fatalf("Mistaken identity %d of %v", follower.me, follower.peerAddresses)
			}

			if follower2.peerAddresses[follower2.me] != "127.0.0.1:4002" {
				t.Fatalf("Mistaken identity %d of %v", follower2.me, follower2.peerAddresses)
			}

			if srv.me != GetPrimary(srv.currentView, int32(len(srv.peerAddresses))) {
				t.Fatalf("Leader does not think she is the primary")
			}

			if follower.me == GetPrimary(follower.currentView, int32(len(follower.peerAddresses))) {
				t.Fatalf("Follower thinks she is primary")
			}

			if follower2.me == GetPrimary(follower2.currentView, int32(len(follower2.peerAddresses))) {
				t.Fatalf("Follower thinks she is primary")
			}

			if int(follower.commitIndex) != len(follower.CommitChan) {
				t.Fatalf("Follower failed to apply committed indices.  Commit channel is of length %d, should be %d", len(follower.CommitChan), follower.commitIndex)
			}

			if int(srv.commitIndex) != len(srv.CommitChan) {
				t.Fatalf("Leader failed to apply committed indices.  Commit channel is of length %d, should be %d", len(srv.CommitChan), srv.commitIndex)
			}

			// log.Printf("Number of commits backed up in leader is %d", len(srv.CommitChan))
			log.Printf("Leader: %v\n1: %v\n2: %v", srv.peerAddresses, follower.peerAddresses, follower2.peerAddresses)
		})
	}
}
