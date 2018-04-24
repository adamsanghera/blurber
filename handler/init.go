package handler

import (
	"log"
	"net"
	"os"

	blurb "github.com/adamsanghera/blurber-protobufs/dist/blurb"
	sub "github.com/adamsanghera/blurber-protobufs/dist/subscription"
	user "github.com/adamsanghera/blurber-protobufs/dist/user"
	"google.golang.org/grpc"
)

func getAddress(envName string) string {
	host := os.Getenv(envName + "_HOST")
	port := os.Getenv(envName + "_PORT")

	if host == "" {
		panic("No hostname set")
	}
	if port == "" {
		panic("No port number set")
	}

	ret, err := net.LookupHost(host)
	if err != nil {
		panic(err)
	}

	return ret[0] + ":" + port
}

func init() {
	// Derive gRPC addresses
	blurbAddr := getAddress("BLURB")
	subAddr := getAddress("SUB")
	userAddr := getAddress("REGISTRATION")
	log.Printf("SERVER: Derived address for blurb db: (%s)", blurbAddr)
	log.Printf("SERVER: Derived address for sub db: (%s)", subAddr)
	log.Printf("SERVER: Derived address for user db: (%s)", userAddr)

	// Connect to gRPC services

	// Blurb database
	connBlurb, err := grpc.Dial(blurbAddr, grpc.WithInsecure())
	if err != nil {
		log.Printf("SERVER: Failed to connect to BlurbDB at %s", blurbAddr)
	}
	blurbDB = blurb.NewBlurbDBClient(connBlurb)
	log.Printf("SERVER: Successfully created a connection to blurb db")

	// Subscription database
	connSub, err := grpc.Dial(subAddr, grpc.WithInsecure())
	if err != nil {
		log.Printf("SERVER: Failed to connect to SubDB at %s", subAddr)
	}
	subDB = sub.NewSubscriptionDBClient(connSub)
	log.Printf("SERVER: Successfully created a connection to sub db")

	// User database
	userSub, err := grpc.Dial(userAddr, grpc.WithInsecure())
	if err != nil {
		log.Printf("SERVER: Failed to connect to UserDB at %s", userAddr)
	}
	userDB = user.NewUserDBClient(userSub)
	log.Printf("SERVER: Successfully created a connection to user db")

	log.Printf("SERVER: Initialized! All gRPC connections secured\n")
}
