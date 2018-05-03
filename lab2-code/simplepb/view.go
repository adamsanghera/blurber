package simplepb

import "github.com/adamsanghera/blurber-protobufs/dist/replication"

// ViewChange is the RPC handler to process ViewChange RPC.
func (srv *PBServer) ViewChange(args *replication.VCArgs, reply *replication.VCReply) {
	// Your code here
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if args.View > srv.currentView {
		srv.currentView = args.View
		srv.status = VIEWCHANGE
		reply.Success = true
		reply.Log = srv.log
		reply.LastNormalView = srv.lastNormalView
	} else {
		reply.Success = false
	}
}

// StartView is the RPC handler to process StartView RPC.
func (srv *PBServer) StartView(args *replication.SVArgs, reply *replication.SVReply) {
	// Your code here
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if srv.currentView < args.View {
		srv.currentView = args.View
		srv.log = args.Log
		srv.status = NORMAL
		srv.lastNormalView = srv.currentView
	}
}

// determineNewViewLog is invoked to determine the log for the newView based on
// the collection of replies for successful ViewChange requests.
// if a quorum of successful replies exist, then ok is set to true.
// otherwise, ok = false.
func (srv *PBServer) determineNewViewLog(successReplies []*replication.VCReply) (
	ok bool, newViewLog []*replication.Command) {
	// Your code here
	newViewLog = srv.log
	minView := srv.lastNormalView
	for _, v := range successReplies {
		if v.LastNormalView < minView {
			newViewLog = v.Log
			minView = v.LastNormalView
		}
	}
	return ok, newViewLog
}
