package subscription

import (
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"golang.org/x/net/context"

	"github.com/adamsanghera/blurber-protobufs/dist/common"
	"github.com/adamsanghera/blurber-protobufs/dist/replication"
	subpb "github.com/adamsanghera/blurber-protobufs/dist/subscription"

	rep "github.com/adamsanghera/blurber/replication-daemon"
)

type LedgerServer struct {
	ledger                   *LocalLedger
	addr                     string
	replicationAdddress      string
	replicationLeaderAddress string

	replicationDaemon *rep.PBServer
}

func (ls *LedgerServer) ProcessCommands() {
	log.Printf("SUB-SERVER: Processing Command loop begins..")

	for cmd := range ls.replicationDaemon.CommitChan {
		log.Printf("Received cmd %v", cmd)

		switch cmd.Cmd.TypeUrl {
		case "sub/add/" + proto.MessageName(&subpb.Subscription{}):
			var sub subpb.Subscription
			err := ptypes.UnmarshalAny(cmd.Cmd, &sub)
			if err != nil {
				panic("Failed to unmarshal command")
			}

			ls.ledger.AddSub(sub.Follower.UserID, sub.Leader.UserID)
			break
		case "sub/delete/" + proto.MessageName(&subpb.Subscription{}):
			var sub subpb.Subscription
			err := ptypes.UnmarshalAny(cmd.Cmd, &sub)
			if err != nil {
				panic("Failed to unmarshal command")
			}

			ls.ledger.RemoveSub(sub.Follower.UserID, sub.Leader.UserID)
			break
		case "sub/delete/" + proto.MessageName(&common.UserID{}):
			var uid common.UserID
			err := ptypes.UnmarshalAny(cmd.Cmd, &uid)
			if err != nil {
				panic("Failed to unmarshal command")
			}

			ls.ledger.RemoveUser(uid.UserID)
			break
		}
	}
}

// NewLedgerServer creates a new ledger server, supported by a command replication daemon
func NewLedgerServer(addr string, replicationAddress string, replicationLeaderAddress string) *LedgerServer {
	// Register proto types
	proto.RegisterType((*subpb.Subscription)(nil), "subpb.Subscription")
	proto.RegisterType((*common.UserID)(nil), "common.UserID")

	// Init the ledger server and daemon
	ls := &LedgerServer{
		ledger:                   NewLocalLedger(),
		addr:                     addr,
		replicationAdddress:      replicationAddress,
		replicationLeaderAddress: replicationLeaderAddress,
		replicationDaemon:        rep.NewReplicationDaemon(replicationAddress, replicationLeaderAddress),
	}

	// Spawn command processor routine, that feeds off the Replication Daemon
	go ls.ProcessCommands()
	return ls
}

func (ls *LedgerServer) Add(ctx context.Context, in *subpb.Subscription) (*common.Empty, error) {
	serialized, err := proto.Marshal(in)
	if err != nil {
		panic(err)
	}

	cmd := &replication.Command{
		Cmd: &any.Any{
			TypeUrl: "sub/add/" + proto.MessageName(in),
			Value:   serialized,
		},
	}

	return &common.Empty{}, ls.replicationDaemon.Replicate(cmd)
}

func (ls *LedgerServer) Delete(ctx context.Context, in *subpb.Subscription) (*common.Empty, error) {
	serialized, err := proto.Marshal(in)
	if err != nil {
		panic(err)
	}

	cmd := &replication.Command{
		Cmd: &any.Any{
			TypeUrl: "sub/delete/" + proto.MessageName(in),
			Value:   serialized,
		},
	}

	return &common.Empty{}, ls.replicationDaemon.Replicate(cmd)
}

func (ls *LedgerServer) DeletePresenceOf(ctx context.Context, in *common.UserID) (*common.Empty, error) {
	serialized, err := proto.Marshal(in)
	if err != nil {
		panic(err)
	}

	cmd := &replication.Command{
		Cmd: &any.Any{
			TypeUrl: "sub/delete/" + proto.MessageName(in),
			Value:   serialized,
		},
	}

	return &common.Empty{}, ls.replicationDaemon.Replicate(cmd)
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
