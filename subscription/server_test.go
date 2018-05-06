package subscription

import (
	"context"
	"testing"
	"time"

	"github.com/adamsanghera/blurber-protobufs/dist/common"
	subpb "github.com/adamsanghera/blurber-protobufs/dist/subscription"
)

func TestLedgerServer_Add(t *testing.T) {
	type args struct {
		in []*subpb.Subscription
	}
	tests := []struct {
		name string
		args args
	}{
		struct {
			name string
			args args
		}{
			name: "yea",
			args: args{
				in: []*subpb.Subscription{
					&subpb.Subscription{
						Follower: &common.UserID{
							UserID: 1,
						},
						Leader: &common.UserID{
							UserID: 2,
						},
					},
					&subpb.Subscription{
						Follower: &common.UserID{
							UserID: 2,
						},
						Leader: &common.UserID{
							UserID: 3,
						},
					},
					&subpb.Subscription{
						Follower: &common.UserID{
							UserID: 3,
						},
						Leader: &common.UserID{
							UserID: 4,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := NewLedgerServer("127.0.0.1:4000", "127.0.0.1:4001")

			for _, sub := range tt.args.in {
				ls.Add(context.Background(), sub)
			}

			for _, sub := range tt.args.in {
				ls.Delete(context.Background(), sub)
			}

			for _, sub := range tt.args.in {
				ls.Add(context.Background(), sub)
			}

			ls.DeletePresenceOf(context.Background(), &common.UserID{
				UserID: 3,
			})

			time.Sleep(1 * time.Second)

			if ls, _ := ls.ledger.GetFollowers(3); len(ls) != 0 {
				t.Fatalf("User 3 followers exist after being deleted")
			}

			if ls, _ := ls.ledger.GetLeaders(2); len(ls) != 0 {
				t.Fatalf("User 3 leaders exist after being deleted %v", ls)
			}
		})
	}
}
