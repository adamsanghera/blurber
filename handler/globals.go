package handler

import (
	blurb "github.com/adamsanghera/blurber/protobufs/dist/blurb"
	sub "github.com/adamsanghera/blurber/protobufs/dist/subscription"
	userpb "github.com/adamsanghera/blurber/protobufs/dist/user"
)

var userDB userpb.UserDBClient
var blurbDB blurb.BlurbDBClient
var subDB sub.SubscriptionDBClient
