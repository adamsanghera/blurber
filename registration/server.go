package registration

import (
	"golang.org/x/net/context"

	"github.com/adamsanghera/blurber-protobufs/dist/common"
	userpb "github.com/adamsanghera/blurber-protobufs/dist/user"
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

func (ls *LedgerServer) Add(ctx context.Context, in *userpb.Credentials) (*common.Empty, error) {
	return &common.Empty{}, ls.ledger.AddNewUser(in.Username, in.Password)
}

func (ls *LedgerServer) Delete(ctx context.Context, in *userpb.Username) (*common.Empty, error) {
	return &common.Empty{}, ls.ledger.Remove(in.Username)
}

func (ls *LedgerServer) GetID(ctx context.Context, in *userpb.Username) (*common.UserID, error) {
	name, err := ls.ledger.GetUserID(in.Username)
	return &common.UserID{
		UserID: name,
	}, err
}

func (ls *LedgerServer) LogIn(ctx context.Context, in *userpb.Credentials) (*userpb.Token, error) {
	token, err := ls.ledger.LogIn(in.Username, in.Password)
	return &userpb.Token{Token: token}, err
}

func (ls *LedgerServer) CheckIn(ctx context.Context, in *userpb.SessionCredentials) (*userpb.Token, error) {
	token, err := ls.ledger.CheckIn(in.Username, in.Token)
	return &userpb.Token{Token: token}, err
}

func (ls *LedgerServer) CheckOut(ctx context.Context, in *userpb.SessionCredentials) (*common.Empty, error) {
	err := ls.ledger.CheckOut(in.Username, in.Token)
	return &common.Empty{}, err
}
