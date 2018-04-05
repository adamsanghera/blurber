package handler

import (
	"os"
)

func init() {
	// Dev only
	if os.Getenv("DEBUG") == "1" {
		userDB.AddNewUser("dev", "root")
		userDB.AddNewUser("adam", "root")
		userDB.AddNewUser("hieu", "root")

		blurbDB.AddNewBlurb(0, "hello", "dev")
		blurbDB.AddNewBlurb(1, "world", "adam")
		blurbDB.AddNewBlurb(2, "wassup", "hieu")
		blurbDB.AddNewBlurb(2, "hieu's day is good", "hieu")
		blurbDB.AddNewBlurb(2, "another day in life of hieu", "hieu")
	}
}
