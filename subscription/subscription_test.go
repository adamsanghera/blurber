package subscription

import (
	"testing"
)

func TestLocalLedger_GetLeaders(t *testing.T) {
	tests := []struct {
		name     string
		leaders  []int32
		follower int32
	}{
		{
			"1",
			[]int32{1389, 8492, 48297, 5774, 48},
			314,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ll := NewLocalLedger()

			for _, leader := range tt.leaders {
				ll.AddSub(tt.follower, leader)
			}

			ll.RemoveUser(tt.leaders[0])
			ll.RemoveUser(tt.follower)

			ret, err := ll.GetLeaders(tt.follower)
			if err != nil {
				t.Fatalf("Couldn't retrieve leaders for %d", tt.follower)
			}

			if len(ret) != 0 {
				t.Fatalf("Wrong number of leaders: %d", len(ret))
			}
		})
	}
}
