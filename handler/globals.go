package handler

import (
	blurb "github.com/adamsanghera/blurber/protobufs/dist/blurb"
	sub "github.com/adamsanghera/blurber/protobufs/dist/subscription"
	reg "github.com/adamsanghera/blurber/registration"
)

var userDB = reg.NewLocalLedger()

var blurbDB blurb.BlurbDBClient
var subDB sub.SubscriptionDBClient
