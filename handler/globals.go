package handler

import (
	"../blurb"
	reg "../registration"
)

var userDB = reg.NewLocalLedger()
var lbl = blurb.NewLocalLedger()
