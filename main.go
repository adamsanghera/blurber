package main

import (
	auth "github.com/adamsanghera/authenticator"
	"github.com/adamsanghera/badge"
)

type LocalLogin struct {
	auth   auth.Authenticator
	minter badge.Minter
}

func main() {
	a, err := auth.NewLocalPasswordAuth(0)
	if err != nil {
		panic(err)
	}

	l := LocalLogin{
		auth:   a,
		minter: badge.NewRandomTokenMinter(256),
	}
}
