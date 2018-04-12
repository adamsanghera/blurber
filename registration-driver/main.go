package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc/reflection"

	userpb "github.com/adamsanghera/blurber/protobufs/dist/user"
	"github.com/adamsanghera/blurber/registration"

	"google.golang.org/grpc"
)

func createListenAddr() string {
	listenPort := os.Getenv("REGISTRATION_PORT")

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

	log.Printf("RegistrationServer: Derived address: (%s)", addr)

	srv := registration.NewLedgerServer(addr)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	userpb.RegisterUserDBServer(s, srv)

	log.Printf("RegistrationServer: Starting on (%s)", addr)

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s\n", err.Error())
	}
}
