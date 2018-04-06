package handler

import (
	"github.com/adamsanghera/blurber/blurb"
	reg "github.com/adamsanghera/blurber/registration"
	sub "github.com/adamsanghera/blurber/subscription"
)

var userDB = reg.NewLocalLedger()
var blurbDB = blurb.NewLocalLedger()
var subDB = sub.NewLocalLedger()
