package handler

import (
	"../blurb"
	reg "../registration"
	sub "../subscription"
)

var userDB = reg.NewLocalLedger()
var blurbDB = blurb.NewLocalLedger()
var subDB = sub.NewLocalLedger()
