package handler

import (
	"log"
	"net"
	"os"

	blurb "github.com/adamsanghera/blurber/protobufs/dist/blurb"
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
	log.Printf("Derived address for blurb db: (%s)", blurbAddr)

	// Connect to gRPC addresses
	conn, err := grpc.Dial(blurbAddr, grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect to BlurbDB at %s", blurbAddr)
	}
	blurbDB = blurb.NewBlurbDBClient(conn)

	log.Printf("Successfully created a connection to blurb db")

	// Dev only
	if os.Getenv("DEBUG") == "1" {
		userDB.AddNewUser("dev", "root")
		userDB.AddNewUser("adam", "root")
		userDB.AddNewUser("hieu", "root")

		// blurbDB.AddNewBlurb(0, "hello", "dev")
		// blurbDB.AddNewBlurb(1, "world", "adam")
		// blurbDB.AddNewBlurb(2, "wassup", "hieu")
		// blurbDB.AddNewBlurb(2, "hieu's day is good", "hieu")
		// blurbDB.AddNewBlurb(2, "another day in life of hieu", "hieu")
	}
}
