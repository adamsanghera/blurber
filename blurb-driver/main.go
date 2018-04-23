package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc/reflection"

	blurbpb "github.com/adamsanghera/blurber-protobufs/dist/blurb"
	"github.com/adamsanghera/blurber/blurb"

	"google.golang.org/grpc"
)

func createListenAddr() string {
	listenPort := os.Getenv("BLURB_PORT")

	// if listenHost == "" {
	// 	panic("No hostname set")
	// }
	if listenPort == "" {
		panic("No port number set")
	}

	return ":" + listenPort
}

func main() {
	addr := createListenAddr()

	log.Printf("BlurbServer: Derived address: (%s)", addr)

	srv := blurb.NewLedgerServer(addr)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	blurbpb.RegisterBlurbDBServer(s, srv)

	log.Printf("BlurbServer: Starting on (%s)", addr)

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s\n", err.Error())
	}
}
