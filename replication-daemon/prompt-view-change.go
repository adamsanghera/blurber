package pbdaemon

import (
	"context"
	"log"

	"github.com/adamsanghera/blurber-protobufs/dist/replication"
)

// PromptViewChange inits the view change protocol to move to the newView
// It does not block waiting for the view change process to complete.
func (srv *PBServer) PromptViewChange(newView int32) {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	log.Printf("Received PVC request with view num %d", newView)

	newPrimary := GetPrimary(newView, int32(len(srv.peers)))

	if newPrimary != srv.me { //only primary of newView should do view change
		return
	} else if newView <= srv.currentView {
		return
	}
	vcArgs := &replication.VCArgs{
		View: newView,
	}
	vcReplyChan := make(chan *replication.VCReply, len(srv.peers))
	// send ViewChange to all servers including myself
	for i := 0; i < len(srv.peers); i++ {
		go func(server int) {
			reply, err := srv.peers[server].ViewChange(context.Background(), vcArgs)
			if err == nil {
				vcReplyChan <- reply
			} else {
				vcReplyChan <- nil
			}
		}(i)
	}

	// wait to receive ViewChange replies
	// if view change succeeds, send StartView RPC
	go func() {
		var successReplies []*replication.VCReply
		var nReplies int
		majority := len(srv.peers)/2 + 1
		for r := range vcReplyChan {
			nReplies++
			if r != nil && r.Success {
				successReplies = append(successReplies, r)
			}
			if nReplies == len(srv.peers) || len(successReplies) == majority {
				log.Printf("Majority of respondents said yes to new view %d", newView)
				break
			}
		}
		ok, newLog := srv.determineNewViewLog(successReplies)
		if !ok {
			return
		}
		log.Printf("New log determined for view %d", newView)
		svArgs := &replication.SVArgs{
			View: vcArgs.View,
			Log:  newLog,
		}
		// send StartView to all servers including myself
		for i := 0; i < len(srv.peers); i++ {
			go func(server int) {
				log.Printf("Sending startview to %d for view %d", server, newView)
				srv.peers[server].StartView(context.Background(), svArgs)
			}(i)
		}
	}()
}
