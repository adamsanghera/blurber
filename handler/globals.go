package handler

import (
	"../blurb"
	reg "../registration"
)

var userDB = reg.NewLocalLedger()
var increasingCounter int

var lbl = blurb.NewLocalLedger()
