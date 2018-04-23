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
	log.Printf("Derived address for blurb db: (%s)", blurbAddr)
	log.Printf("Derived address for sub db: (%s)", subAddr)
	log.Printf("Derived address for user db: (%s)", userAddr)

	// Connect to gRPC addresses
	connBlurb, err := grpc.Dial(blurbAddr, grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect to BlurbDB at %s", blurbAddr)
	}
	blurbDB = blurb.NewBlurbDBClient(connBlurb)
	log.Printf("Successfully created a connection to blurb db")

	connSub, err := grpc.Dial(subAddr, grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect to SubDB at %s", subAddr)
	}
	subDB = sub.NewSubscriptionDBClient(connSub)
	log.Printf("Successfully created a connection to sub db")

	userSub, err := grpc.Dial(userAddr, grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect to UserDB at %s", userAddr)
	}
	userDB = user.NewUserDBClient(userSub)
	log.Printf("Successfully created a connection to user db")

	// Dev only
	if os.Getenv("DEBUG") == "1" {
		// userDB.AddNewUser("dev", "root")
		// userDB.AddNewUser("adam", "root")
		// userDB.AddNewUser("hieu", "root")

		// blurbDB.AddNewBlurb(0, "hello", "dev")
		// blurbDB.AddNewBlurb(1, "world", "adam")
		// blurbDB.AddNewBlurb(2, "wassup", "hieu")
		// blurbDB.AddNewBlurb(2, "hieu's day is good", "hieu")
		// blurbDB.AddNewBlurb(2, "another day in life of hieu", "hieu")
	}
}
