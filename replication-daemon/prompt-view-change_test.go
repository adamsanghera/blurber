package pbdaemon

import (
	"log"
	"sync"
	"testing"
	"time"

	"github.com/adamsanghera/blurber-protobufs/dist/replication"
	"github.com/golang/protobuf/ptypes/any"
)

func TestPBServer_PromptViewChange(t *testing.T) {
	type fields struct {
		mu             *sync.Mutex
		peers          []replication.ReplicationClient
		peerAddresses  []string
		me             int32
		currentView    int32
		status         int32
		lastNormalView int32
		log            []*replication.Command
		commitIndex    int32
		CommitChan     chan *replication.Command
		prepChan       chan *callbackArg
	}
	type args struct {
		// newView int32

		leaderAddr         string
		firstFollowerAddr  string
		secondFollowerAddr string
	}
	tests := []struct {
		name string
		// fields fields
		args args
	}{
		// TODO: Add test cases.
		struct {
			name string
			args args
		}{
			"test1",
			args{
				leaderAddr:         "127.0.0.1:4000",
				firstFollowerAddr:  "127.0.0.1:4001",
				secondFollowerAddr: "127.0.0.1:4002",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Assuming NewReplicationDaemon is already tested

			// Spawn 3 daemons
			srv := NewReplicationDaemon(tt.args.leaderAddr, tt.args.leaderAddr)
			firstFollower := NewReplicationDaemon(tt.args.firstFollowerAddr, tt.args.leaderAddr)
			secondFollower := NewReplicationDaemon(tt.args.secondFollowerAddr, tt.args.leaderAddr)

			// Replicate 1 command first
			srv.Replicate(&replication.Command{
				Cmd: &any.Any{
					TypeUrl: "Add",
					Value:   []byte("hieu"),
				},
			})

			time.Sleep(2 * time.Second)

			if len(firstFollower.log) != len(srv.log) || len(secondFollower.log) != len(srv.log) {
				t.Fatalf("Failed to replicate log to follower")
				log.Printf("{%v}", firstFollower.log)
				log.Printf("{%v}", secondFollower.log)
				log.Printf("{%v}", srv.log)
			}

			// Crash the oldLeader
			srv.StopGRPCServer()

			time.Sleep(1 * time.Second)

			// Replicate command on oldLeader
			srv.Replicate(&replication.Command{
				Cmd: &any.Any{
					TypeUrl: "Add",
					Value:   []byte("mich"),
				},
			})

			time.Sleep(1 * time.Second)

			if len(firstFollower.log) == len(srv.log) || len(secondFollower.log) == len(srv.log) {
				log.Printf("{%v}", firstFollower.log)
				log.Printf("{%v}", secondFollower.log)
				log.Printf("{%v}", srv.log)
				t.Fatalf("Failed to disconnect oldLeader")
			}

			// Change view

			// srv.PromptViewChange(tt.args.newView)

			// Get new leader

			// Replicate commands at newPrimary

			// Try to replicate commands at oldPrimary

			// Reconnect oldPrimary

			// Replicate commands at newPrimary

		})
	}
}
