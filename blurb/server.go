package blurb

import (
	"golang.org/x/net/context"

	"github.com/adamsanghera/blurber/protobufs/dist/blurb"
	"github.com/adamsanghera/blurber/protobufs/dist/common"
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

func (ls *LedgerServer) Add(c context.Context, b *blurb.NewBlurb) (*common.Empty, error) {
	ls.ledger.AddNewBlurb(int(b.Author.UserID), b.Content, b.Username)
	return &common.Empty{}, nil
}

func (ls *LedgerServer) Delete(c context.Context, bi *blurb.BlurbIndex) (*common.Empty, error) {
	ls.ledger.RemoveBlurb(int(bi.Author.UserID), int(bi.BlurbID))
	return &common.Empty{}, nil
}

func (ls *LedgerServer) GenerateFeed(c context.Context, fp *blurb.FeedParameters) (*blurb.Blurbs, error) {
	leaders := make([]int, len(fp.LeaderIDs))
	for k, v := range fp.LeaderIDs {
		leaders[k] = int(v.UserID)
	}
	blurbs := ls.ledger.GenerateFeed(int(fp.RequestorID.UserID), leaders)

	ret := make([]*blurb.Blurb, len(blurbs))

	for k, v := range blurbs {
		ret[k] = &blurb.Blurb{
			BlurbID:     int32(v.BID),
			Content:     v.Content,
			CreatorName: v.CreatorName,
			Timestamp:   v.Timestamp,
			UnixTime:    v.Time.Unix(),
		}
	}
	return &blurb.Blurbs{
		Blurbs: ret,
	}, nil
}

func (ls *LedgerServer) GetRecentBy(c context.Context, uid *common.UserID) (*blurb.Blurbs, error) {
	blurbs := ls.ledger.GetRecentBlurbsBy(int(uid.UserID))

	ret := make([]*blurb.Blurb, len(blurbs))

	for k, v := range blurbs {
		ret[k] = &blurb.Blurb{
			BlurbID:     int32(v.BID),
			Content:     v.Content,
			CreatorName: v.CreatorName,
			Timestamp:   v.Timestamp,
			UnixTime:    v.Time.Unix(),
		}
	}

	return &blurb.Blurbs{
		Blurbs: ret,
	}, nil
}

func (ls *LedgerServer) DeleteHistoryOf(c context.Context, uid *common.UserID) (*common.Empty, error) {
	ls.ledger.RemoveAllBlurbsBy(int(uid.UserID))
	return &common.Empty{}, nil
}

func (ls *LedgerServer) InvalidateFeedCache(c context.Context, uid *common.UserID) (*common.Empty, error) {
	ls.ledger.InvalidateCache(int(uid.UserID))
	return &common.Empty{}, nil
}
