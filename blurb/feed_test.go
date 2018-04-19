package blurb

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestLocalLedger_GenerateFeed(t *testing.T) {
	tests := []struct {
		name      string
		numBlurbs int32
		sources   []int32
	}{}

	for numFriends := 0; numFriends < 10; numFriends++ {
		// Make the friend list
		friendList := make([]int32, numFriends)
		for idx := 1; idx <= numFriends; idx++ {
			friendList[idx-1] = int32(idx)
		}

		// Generate the test cases for each friend
		for testSize := 0; testSize < 20; testSize++ {
			tests = append(tests, struct {
				name      string
				numBlurbs int32
				sources   []int32
			}{
				fmt.Sprintf("Testing %d friends with uniform blurb lists of %d each", numFriends, testSize),
				int32(testSize),
				friendList,
			})
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ll := NewLocalLedger()

			for friendIdx := 0; friendIdx < len(tt.sources); friendIdx++ {
				for idx := 0; idx < int(tt.numBlurbs); idx++ {
					ll.AddNewBlurb(tt.sources[friendIdx], "I am the best"+strconv.Itoa(idx), "Adam clone number "+strconv.Itoa(int(tt.sources[friendIdx])))
				}
			}

			feed := ll.GenerateFeed(0, tt.sources)
			if len(feed) != min(min(int(tt.numBlurbs)*len(tt.sources), len(tt.sources)*10), 100) {
				t.Fatalf("Feed is of the wrong length (%d) when it should be (%d)", len(feed), min(int(tt.numBlurbs)*len(tt.sources), len(tt.sources)*10))
				t.Fail()
			}

			for idx := range feed {
				if idx > 1 {
					iTime := time.Unix(feed[idx].UnixTime, 0)
					jTime := time.Unix(feed[idx-1].UnixTime, 0)
					if iTime.After(jTime) {
						t.Errorf("Feed is out of order: {%s} is placed before {%s}", feed[idx-1].Timestamp, feed[idx].Timestamp)
					}
				}
			}
		})
	}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
