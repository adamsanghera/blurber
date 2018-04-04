package blurb

import (
	"fmt"
	"strconv"
	"testing"
)

func TestLocalLedger_RemoveAllBlurbsBy(t *testing.T) {
	tests := []struct {
		name      string
		numBlurbs int
	}{
		{"Adding blurbs less than cache size", 10},
		{"Adding blurbs greater than cache size", 20},
		{"No-op", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ll := NewLocalLedger()
			for idx := 0; idx < tt.numBlurbs; idx++ {
				ll.AddNewBlurb(0, "i m d best "+strconv.Itoa(idx), "adam")
			}
			ll.RemoveAllBlurbsBy(0)

			if len(ll.GetBlurbsCreatedBy(0)) > 0 {
				t.Errorf("GetBlurbsCreatedBy returned a non-empty (size %d) list", len(ll.GetBlurbsCreatedBy(0)))
			}
			if len(ll.GetRecentBlurbsBy(0)) > 0 {
				t.Errorf("GetRecentBlurbsBy returned a non-empty (size %d) list", len(ll.GetRecentBlurbsBy(0)))
			}
		})
	}
}

func TestLocalLedger_RemoveBlurb(t *testing.T) {
	tests := []struct {
		name              string
		numBlurbsExisting int
		blurbToDelete     int
	}{}

	bDel := 0
	numBlurbsExisting := []int{10, 15, 30}

	for _, numExisting := range numBlurbsExisting {
		for bDel < numExisting {
			tests = append(tests, struct {
				name              string
				numBlurbsExisting int
				blurbToDelete     int
			}{
				fmt.Sprintf("Deleting %d from a history of length %d", bDel, numExisting),
				numExisting,
				bDel,
			})
			bDel++
		}
		bDel = 0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ll := NewLocalLedger()
			for idx := 0; idx < tt.numBlurbsExisting; idx++ {
				ll.AddNewBlurb(0, "i m d best "+strconv.Itoa(idx), "adam")
			}

			ll.RemoveBlurb(0, tt.blurbToDelete)

			allBlurbs := ll.GetBlurbsCreatedBy(0)
			cachedBlurbs := ll.GetRecentBlurbsBy(0)

			for _, b := range allBlurbs {
				if b.BID == tt.blurbToDelete {
					t.Errorf("History contains the supposedly deleted blurb")
				}
			}

			for _, b := range cachedBlurbs {
				if b.BID == tt.blurbToDelete {
					t.Errorf("Cache contains the supposedly deleted blurb")
				}
			}
		})
	}
}
