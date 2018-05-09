package pbdaemon

//
// This is a outline of primary-backup replication based on a simplifed version of Viewstamp replication.
//
//
//

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/adamsanghera/blurber-protobufs/dist/replication"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// the 3 possible server status
const (
	NORMAL = iota
	VIEWCHANGE
	RECOVERING
)

type callbackArg struct {
	callback chan bool
	args     *replication.PrepareArgs
	handled  bool
}

// PBServer defines the state of a replica server (either primary or backup)
type PBServer struct {
	mu             *sync.Mutex                     // Lock to protect shared access to this peer's state
	peers          []replication.ReplicationClient // RPC end points of all peers
	peerAddresses  []string                        // Addresses of all peers
	me             int32                           // this peer's index into peers[]
	currentView    int32                           // what this peer believes to be the current active view
	status         int32                           // the server's current status (NORMAL, VIEWCHANGE or RECOVERING)
	lastNormalView int32                           // the latest view which had a NORMAL status

	log         []*replication.Command // the log of "commands"
	commitIndex int32                  // all log entries <= commitIndex are considered to have been committed.
	CommitChan  chan *replication.Command

	prepChan chan *callbackArg // Channel used by prep calls to communicate with the central prep-processor

	grpcServer *grpc.Server
}

// GetPrimary is an auxilary function that returns the server index of the
// primary server given the view number (and the total number of replica servers)
func GetPrimary(view int32, nservers int32) int32 {
	return view % nservers
}

func initGRPCServer(thisAddress string, srv *PBServer) {
	log.Printf("ReplicationD: Registering to listen on (%s)", thisAddress)
	lis, err := net.Listen("tcp", thisAddress)
	if err != nil {
		panic(err)
	}
	srv.grpcServer = grpc.NewServer()
	replication.RegisterReplicationServer(srv.grpcServer, srv)
	reflection.Register(srv.grpcServer)

	log.Printf("Replicatio nD: Registered successfully (%s)", thisAddress)

	if err = srv.grpcServer.Serve(lis); err != nil {
		panic("ReplicationD: Failed to start gRPC server")
	}
}

// Home grown find method for slice since go standard lib doesn't have one
// Maybe this is why: https://github.com/golang/go/wiki/InterfaceSlice
func findSlice(addr string, allAddrs []string) int {
	for i, e := range allAddrs {
		if e == addr {
			return i
		}
	}
	return -1

}

func (srv *PBServer) disconnectPeer(addr string) error {
	// do we really need to disconnect client conn?

	// Get index i of the peer addr to be disconnected
	i := findSlice(addr, srv.peerAddresses)
	srv.peers = append(srv.peers[:i], srv.peers[i+1:]...)
	srv.peerAddresses = append(srv.peerAddresses[:i], srv.peerAddresses[i+1:]...)

	return nil
}

// StopGRPCServer is used in unit testing to simulate server failure
func (srv *PBServer) StopGRPCServer() {
	log.Printf("Stopping gRPC server %v", srv.me)

	// Disconnect peers
	for i, addr := range srv.peerAddresses {
		if int32(i) != srv.me {
			srv.disconnectPeer(addr)
		}
	}

	srv.grpcServer.Stop()
}

func (srv *PBServer) PropagateCommitsUnsafe() {
	for idx, client := range srv.peers {
		if idx != int(srv.me) {
			go client.Prepare(context.Background(), &replication.PrepareArgs{View: srv.currentView, PrimaryCommit: srv.commitIndex, Index: -1, Entry: nil})
		}
	}
}

// NewReplicationDaemon spawns a new replication daemon.
// If thisAddress and leaderAddress are not the same,
// then the new daemon will start in recovery mode.
func NewReplicationDaemon(thisAddress string, leaderAddress string) *PBServer {
	srv := &PBServer{
		mu:             &sync.Mutex{},
		peers:          make([]replication.ReplicationClient, 0),
		peerAddresses:  make([]string, 0),
		me:             0,
		currentView:    0,
		status:         NORMAL,
		lastNormalView: 0,
		log:            make([]*replication.Command, 0),
		commitIndex:    0,
		CommitChan:     make(chan *replication.Command, 1000),
		prepChan:       make(chan *callbackArg),
	}

	srv.mu.Lock()
	defer srv.mu.Unlock()

	log.Printf("New Daemon created at %s, pointing at %s", thisAddress, leaderAddress)

	// Init log
	srv.log = append(srv.log, &replication.Command{})

	// Initting daemon connections
	if thisAddress != leaderAddress {
		log.Printf("ReplicationD: Spawning as follower")
		srv.me = 1
		err := srv.connectPeerUnsafe(leaderAddress)
		if err != nil {
			panic(err)
		}
		err = srv.connectPeerUnsafe(thisAddress)
		if err != nil {
			panic(err)
		}
	} else {
		log.Printf("ReplicationD: Spawning as leader")
		err := srv.connectPeerUnsafe(thisAddress)
		if err != nil {
			panic(err)
		}
	}

	log.Printf("Spwaning replication processor")
	go srv.prepareProcessor()

	// Registering server
	go initGRPCServer(thisAddress, srv)

	if leaderAddress != thisAddress {
		log.Printf("ReplicationD: Follower entering recovery mode")
		go srv.sendRecovery()
	}

	return srv
}
