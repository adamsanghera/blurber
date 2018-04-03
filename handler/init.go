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

        lbl.AddNewBlurb(0, "hello", "dev")
        lbl.AddNewBlurb(1, "world", "adam")
        lbl.AddNewBlurb(2, "wassup", "hieu")
    }
}
