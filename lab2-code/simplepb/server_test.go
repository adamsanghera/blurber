package simplepb

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/any"

	"github.com/adamsanghera/blurber-protobufs/dist/replication"
	"google.golang.org/grpc"
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

			// Create a replication client
			conn, err := grpc.Dial(tt.args.leaderAddress, grpc.WithInsecure())
			if err != nil {
				t.Fatalf("ReplicationD: Failed to connect to the Replication Daemon at %s", tt.args.leaderAddress)

			}
			client := replication.NewReplicationClient(conn)

			client.Replicate(context.Background(), &replication.Command{
				Cmd: &any.Any{
					TypeUrl: "Add",
					Value:   []byte("adam"),
				},
			})
			client.Replicate(context.Background(), &replication.Command{
				Cmd: &any.Any{
					TypeUrl: "Add",
					Value:   []byte("bob"),
				},
			})
			client.Replicate(context.Background(), &replication.Command{
				Cmd: &any.Any{
					TypeUrl: "Add",
					Value:   []byte("bob"),
				},
			})

			time.Sleep(2 * time.Second)

			if len(srv.log) != 4 {
				t.Fatalf("Failed to store commands")
			}

			follower := NewReplicationDaemon(tt.args.thisAddress, tt.args.leaderAddress)

			time.Sleep(2 * time.Second)

			if len(follower.log) != len(srv.log) {
				t.Fatalf("Failed to replicate log to follower")
			}
		})
	}
}
