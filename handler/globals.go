package handler

import (
	"../blurb"
	reg "../registration"
    sub "../subscription"
)

var userDB = reg.NewLocalLedger()
var lbl = blurb.NewLocalLedger()
var subDB = sub.NewLocalLedger()