package subscription

import (
	"golang.org/x/net/context"

	"github.com/adamsanghera/blurber-protobufs/dist/common"
	subpb "github.com/adamsanghera/blurber-protobufs/dist/subscription"
)

type LedgerServer struct {
	ledger *LocalLedger
	addr   string
}

func NewLedgerServer(addr string) *LedgerServer {
	return &LedgerServer{
		ledger: NewLocalLedger(),
		addr:   addr,
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

func (ls *LedgerServer) GetLeadersOf(ctx context.Context, in *common.UserID) (*subpb.Leaders, error) {
	ret, err := ls.ledger.GetLeaders(in.UserID)

	retList := make([]*common.UserID, len(ret))

	for k := range ret {
		retList[k] = &common.UserID{UserID: ret[k]}
	}

	return &subpb.Leaders{
		Leaders: retList,
	}, err
}
