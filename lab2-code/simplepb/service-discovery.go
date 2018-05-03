package simplepb

// maintainPeerList is a background function, init'd by Start, making use of the directory channel.
func (srv *PBServer) maintainPeerList() {
	for idx := range srv.serviceReportChannel {
		if idx == -1 {
			// Received a recovery message from a new node
			srv.mu.Lock()
			// client := grpc.
		}
	}
}
