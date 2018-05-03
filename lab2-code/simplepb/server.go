package simplepb

//
// This is a outline of primary-backup replication based on a simplifed version of Viewstamp replication.
//
//
//

import (
	"context"
	"sync"

	"github.com/adamsanghera/blurber-protobufs/dist/replication"
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
	mu             sync.Mutex                      // Lock to protect shared access to this peer's state
	peers          []replication.ReplicationClient // RPC end points of all peers
	me             int32                           // this peer's index into peers[]
	currentView    int32                           // what this peer believes to be the current active view
	status         int32                           // the server's current status (NORMAL, VIEWCHANGE or RECOVERING)
	lastNormalView int32                           // the latest view which had a NORMAL status

	log         []*replication.Command // the log of "commands"
	commitIndex int32                  // all log entries <= commitIndex are considered to have been committed.

	prepChan             chan *callbackArg // Channel used by prep calls to communicate with the central prep-processor
	serviceReportChannel chan int          // Channel used by service-discovery routine to maintain peers list.
}

// GetPrimary is an auxilary function that returns the server index of the
// primary server given the view number (and the total number of replica servers)
func GetPrimary(view int32, nservers int32) int32 {
	return view % nservers
}

// IsCommitted is called by tester to check whether an index position
// has been considered committed by this server
func (srv *PBServer) IsCommitted(index int32) (committed bool) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if srv.commitIndex >= index {
		return true
	}
	return false
}

// ViewStatus is called by tester to find out the current view of this server
// and whether this view has a status of NORMAL.
func (srv *PBServer) ViewStatus() (currentView int32, statusIsNormal bool) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	return srv.currentView, srv.status == NORMAL
}

// GetEntryAtIndex is called by tester to return the command replicated at
// a specific log index. If the server's log is shorter than "index", then
// ok = false, otherwise, ok = true
func (srv *PBServer) GetEntryAtIndex(index int) (ok bool, command interface{}) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if len(srv.log) > index {
		return true, srv.log[index]
	}
	return false, command
}

// Make is called by tester to create and initalize a PBServer
// peers is the list of RPC endpoints to every server (including self)
// me is this server's index into peers.
// startingView is the initial view (set to be zero) that all servers start in
func Make(peers []replication.ReplicationClient, me int32, startingView int32) *PBServer {
	srv := &PBServer{
		peers:          peers,
		me:             me,
		currentView:    startingView,
		lastNormalView: startingView,
		status:         NORMAL,
		prepChan:       make(chan *callbackArg),
	}
	// all servers' log are initialized with a dummy command at index 0
	srv.log = append(srv.log, &replication.Command{})

	go srv.prepareProcessor()

	// Your other initialization code here, if there's any
	return srv
}

// exmple code to send an AppendEntries RPC to a server.
// server is the index of the target server in srv.peers[].
// expects RPC arguments in args.
// The RPC library fills in *reply with RPC reply, so caller should pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
func (srv *PBServer) sendPrepare(server int, args *replication.PrepareArgs, reply *replication.PrepareReply) bool {
	reply, err := srv.peers[server].Prepare(context.Background(), args)
	return err == nil
}
