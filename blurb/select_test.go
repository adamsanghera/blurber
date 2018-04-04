package blurb

import (
	"fmt"
	"strconv"
	"testing"
)

func TestLocalLedger_GetBlurbsCreatedBy(t *testing.T) {
	tests := []struct {
		name      string
		numBlurbs int
	}{}

	for testSize := 0; testSize < 100; testSize++ {
		tests = append(tests, struct {
			name      string
			numBlurbs int
		}{
			fmt.Sprintf("Testing selection with blurb list size %d", testSize),
			testSize,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ll := NewLocalLedger()
			for idx := 0; idx < tt.numBlurbs; idx++ {
				ll.AddNewBlurb(0, "i m d best "+strconv.Itoa(idx), "adam")
			}

			if len(ll.GetBlurbsCreatedBy(0)) != tt.numBlurbs {
				t.Fatalf("GetBlurbsCreatedBy returned a list (size %d) of the wrong size", len(ll.GetBlurbsCreatedBy(0)))
			}
		})
	}
}

func TestLocalLedger_GetRecentBlurbsBy(t *testing.T) {
	tests := []struct {
		name      string
		numBlurbs int
	}{}

	for testSize := 0; testSize < 25; testSize++ {
		tests = append(tests, struct {
			name      string
			numBlurbs int
		}{
			fmt.Sprintf("Testing selection with blurb list size %d", testSize),
			testSize,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ll := NewLocalLedger()
			for idx := 0; idx < tt.numBlurbs; idx++ {
				ll.AddNewBlurb(0, "i m d best "+strconv.Itoa(idx), "adam")
			}

			recentBlurbs := ll.GetRecentBlurbsBy(0)

			if len(recentBlurbs) > 10 {
				t.Errorf("GetRecentBlurbsBy returned a list (size %d) of the wrong size", len(recentBlurbs))
			}

			for idx := range recentBlurbs {
				if idx > 1 {
					if recentBlurbs[idx].Time.After(recentBlurbs[idx-1].Time) {
						t.Errorf("GetRecentBlurbsBy returned a list that is out of order: {%s} is placed before {%s}", recentBlurbs[idx-1].Timestamp, recentBlurbs[idx].Timestamp)
					}
				}
			}

		})
	}
}
