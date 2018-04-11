package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc/reflection"

	"github.com/adamsanghera/blurber/blurb"
	blurbpb "github.com/adamsanghera/blurber/protobufs/dist/blurb"

	"google.golang.org/grpc"
)

func createListenAddr() string {
	listenHost := os.Getenv("BLURB_HOST")
	listenPort := os.Getenv("BLURB_PORT")

	// if listenHost == "" {
	// 	panic("No hostname set")
	// }
	if listenPort == "" {
		panic("No port number set")
	}

	ret, err := net.LookupHost(listenHost)
	if err != nil {
		panic(err)
	}

	return ret[0] + ":" + listenPort
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
