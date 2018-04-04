package blurb

import (
	"fmt"
	"strconv"
	"testing"
)

func TestLocalLedger_GenerateFeed(t *testing.T) {
	tests := []struct {
		name      string
		numBlurbs int
		sources   []int
	}{}

	for numFriends := 0; numFriends < 10; numFriends++ {
		// Make the friend list
		friendList := make([]int, numFriends)
		for idx := 1; idx <= numFriends; idx++ {
			friendList[idx-1] = idx
		}

		// Generate the test cases for each friend
		for testSize := 0; testSize < 20; testSize++ {
			tests = append(tests, struct {
				name      string
				numBlurbs int
				sources   []int
			}{
				fmt.Sprintf("Testing %d friends with uniform blurb lists of %d each", numFriends, testSize),
				testSize,
				friendList,
			})
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ll := NewLocalLedger()

			for friendIdx := 0; friendIdx < len(tt.sources); friendIdx++ {
				for idx := 0; idx < tt.numBlurbs; idx++ {
					ll.AddNewBlurb(tt.sources[friendIdx], "I am the best"+strconv.Itoa(idx), "Adam clone number "+strconv.Itoa(tt.sources[friendIdx]))
				}
			}

			feed := ll.GenerateFeed(0, tt.sources)
			if len(feed) != min(min(tt.numBlurbs*len(tt.sources), len(tt.sources)*10), 100) {
				t.Fatalf("Feed is of the wrong length (%d) when it should be (%d)", len(feed), min(tt.numBlurbs*len(tt.sources), len(tt.sources)*10))
				t.Fail()
			}

			for idx := range feed {
				if idx > 1 {
					if feed[idx].Time.After(feed[idx-1].Time) {
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
