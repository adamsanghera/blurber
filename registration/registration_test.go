package registration

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestLocalLedger_LogIn(t *testing.T) {
	type args struct {
		unames []string
		pwds   []string
	}
	tests := []struct {
		name string
		args args
	}{}

	for idx := 0; idx < 10; idx++ {
		testArgs := args{}
		testArgs.unames = make([]string, idx)
		testArgs.pwds = make([]string, idx)
		for jdx := 0; jdx < idx; jdx++ {
			testArgs.unames[jdx] = fmt.Sprintf("name num %d", +jdx)
			testArgs.pwds[jdx] = fmt.Sprintf("pwd num %d", +jdx)
		}
		tests = append(tests, struct {
			name string
			args args
		}{
			fmt.Sprintf("Test with %d concurrent users", idx),
			testArgs,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lul := &LocalLedger{
				userSet:    sync.Map{},
				userID:     sync.Map{},
				pwdMap:     sync.Map{},
				tokenMap:   sync.Map{},
				uidCounter: 0,
			}

			numUsers := len(tt.args.unames)

			for idx := 0; idx < numUsers; idx++ {
				lul.AddNewUser(tt.args.unames[idx], tt.args.pwds[idx])
			}

			for idx := 0; idx < numUsers; idx++ {
				go func(uname string, pwd string) {
					tok, err := lul.LogIn(uname, pwd)
					if err != nil {
						panic(err)
					}
					time.Sleep(time.Millisecond * 1)
					for jdx := 0; jdx < 25; jdx++ {
						tok, err = lul.CheckIn(uname, tok)
						if err != nil {
							panic(err)
						}
						time.Sleep(time.Microsecond * 17)
					}
					lul.CheckOut(uname, pwd)
				}(tt.args.unames[idx], tt.args.pwds[idx])
			}
		})
	}
}
