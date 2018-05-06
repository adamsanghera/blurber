package subscription

import (
	"golang.org/x/net/context"

	"github.com/adamsanghera/blurber-protobufs/dist/common"
	subpb "github.com/adamsanghera/blurber-protobufs/dist/subscription"
	rep "github.com/adamsanghera/blurber/replication-daemon"
)

type LedgerServer struct {
	ledger       *LocalLedger
	addr         string
	internalAddr string

	replicationDaemon *rep.PBServer
}

func NewLedgerServer(addr string, internalAddr string) *LedgerServer {
	return &LedgerServer{
		ledger:       NewLocalLedger(),
		addr:         addr,
		internalAddr: internalAddr,

		replicationDaemon: rep.NewReplicationDaemon(internalAddr, internalAddr),
	}
}

func (ls *LedgerServer) Add(ctx context.Context, in *subpb.Subscription) (*common.Empty, error) {
	ls.ledger.AddSub(in.Follower.UserID, in.Leader.UserID)
	return &common.Empty{}, nil
}

func (ls *LedgerServer) Delete(ctx context.Context, in *subpb.Subscription) (*common.Empty, error) {
	ls.ledger.RemoveSub(in.Follower.UserID, in.Leader.UserID)
	return &common.Empty{}, nil
}

func (ls *LedgerServer) DeletePresenceOf(ctx context.Context, in *common.UserID) (*common.Empty, error) {
	ls.ledger.RemoveUser(in.UserID)
	return &common.Empty{}, nil
}

func (ls *LedgerServer) GetLeadersOf(ctx context.Context, in *common.UserID) (*subpb.Users, error) {
	ret, err := ls.ledger.GetLeaders(in.UserID)

	retList := make([]*common.UserID, len(ret))

	for k := range ret {
		retList[k] = &common.UserID{UserID: ret[k]}
	}

	return &subpb.Users{
		Users: retList,
	}, err
}

func (ls *LedgerServer) GetFollowersOf(ctx context.Context, in *common.UserID) (*subpb.Users, error) {
	ret, err := ls.ledger.GetFollowers(in.UserID)

	retList := make([]*common.UserID, len(ret))

	for k := range ret {
		retList[k] = &common.UserID{UserID: ret[k]}
	}

	return &subpb.Users{
		Users: retList,
	}, err
}
