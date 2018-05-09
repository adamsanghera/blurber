package pbdaemon

import (
	"context"
	"log"

	"github.com/adamsanghera/blurber-protobufs/dist/replication"
)

// ViewChange is the RPC handler to process ViewChange RPC.
func (srv *PBServer) ViewChange(ctx context.Context, args *replication.VCArgs) (*replication.VCReply, error) {
	// Your code here
	srv.mu.Lock()
	defer srv.mu.Unlock()
	log.Printf("Responding to view change request %d", args.View)
	if args.View > srv.currentView {
		log.Printf("Accepting view change request %d", args.View)
		srv.currentView = args.View
		srv.status = VIEWCHANGE
		return &replication.VCReply{
			Success:        true,
			Log:            srv.log,
			LastNormalView: srv.lastNormalView,
		}, nil
	}
	log.Printf("Rejecting view change request %d", args.View)
	return &replication.VCReply{
		Success: false,
	}, nil
}

// StartView is the RPC handler to process StartView RPC.
func (srv *PBServer) StartView(ctx context.Context, args *replication.SVArgs) (*replication.SVReply, error) {
	// Your code here
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if srv.currentView <= args.View {
		log.Printf("Accepting new view %d", args.View)
		srv.currentView = args.View
		srv.log = args.Log
		srv.status = NORMAL
		srv.lastNormalView = srv.currentView
	}

	return &replication.SVReply{}, nil
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
	return true, newViewLog
}
