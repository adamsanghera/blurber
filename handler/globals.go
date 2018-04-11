package handler

import (
	blurb "github.com/adamsanghera/blurber/protobufs/dist/blurb"
	reg "github.com/adamsanghera/blurber/registration"
	sub "github.com/adamsanghera/blurber/subscription"
)

var userDB = reg.NewLocalLedger()

var blurbDB blurb.BlurbDBClient
var subDB = sub.NewLocalLedger()
