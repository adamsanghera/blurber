package handler

import (
	"context"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/adamsanghera/blurber-protobufs/dist/blurb"
	"github.com/adamsanghera/blurber-protobufs/dist/common"
	"github.com/adamsanghera/blurber-protobufs/dist/subscription"
	"github.com/adamsanghera/blurber-protobufs/dist/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type config struct {
	addresses struct {
		blurbDB string
		userDB  string
	}
	clients struct {
		blurbDB blurb.BlurbDBClient
		userDB  user.UserDBClient
		subDBs  map[string]subscription.SubscriptionDBClient
	}
	connections struct {
		blurb *grpc.ClientConn
		user  *grpc.ClientConn
		subs  map[string]*grpc.ClientConn
	}
	leaders struct {
		blurb string // useless
		user  string // useless
		sub   string // useful
	}
	muts struct {
		blurbMut *sync.Mutex
		userMut  *sync.Mutex
		subMut   *sync.Mutex
	}
}

// updateLeaders performs a parallelized poll of the gRPC servers, to determine
// which is the leader
func (c *config) updateLeaders() {
	c.muts.subMut.Lock()
	voteChannel := make(chan *common.ServerInfo)
	for addr, client := range c.clients.subDBs {

		// Spawn a function to get a vote from every client
		go func(vChan chan *common.ServerInfo, addr string, subclient subscription.SubscriptionDBClient) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
			info, err := subclient.Leader(ctx, &common.Empty{})
			if err != nil {
				log.Printf("SERVER CONN: Failed to contact %s, with msg {%s}", addr, err.Error())
			}
			vChan <- info
			cancel()
		}(voteChannel, addr, client)

	}

	voteBox := make(map[string]int)

	for range c.clients.subDBs {
		vote := <-voteChannel
		if _, exists := voteBox[vote.Address]; !exists {
			voteBox[vote.Address] = 0
		}
		voteBox[vote.Address]++
	}

	winner := ""
	winningTally := 0
	for key, value := range voteBox {
		if value > winningTally {
			winner = key
		}
	}

	c.leaders.sub = winner
	c.muts.subMut.Unlock()
}

// maintainConnections calls updateLeaders periodically
func (c *config) maintainConnections() {
	epoch := 0
	for {
		c.updateLeaders()
		log.Printf("SERVER: Leader for epoch %d is %s", epoch, c.leaders.sub)
		time.Sleep(5 * time.Second)
	}
}

func (c *config) toBlurbDB() blurb.BlurbDBClient {
	return c.clients.blurbDB
}

func (c *config) toUserDB() user.UserDBClient {
	return c.clients.userDB
}

func (c *config) toSubDBs() subscription.SubscriptionDBClient {
	// acquire the lock so that we don't concurrently access the subDBs map
	c.muts.subMut.Lock()
	defer c.muts.subMut.Unlock()

	return c.clients.subDBs[c.leaders.sub]
}

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

func gatherAddresses() {
	// Derive gRPC addresses
	configuration.addresses.blurbDB = getAddress("BLURB")
	configuration.addresses.userDB = getAddress("REGISTRATION")
}

func dialServers() {
	dialOpts := grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                5 * time.Second,
		Timeout:             1 * time.Second,
		PermitWithoutStream: true,
	})

	// Blurb database
	var err error
	configuration.connections.blurb, err = grpc.Dial(configuration.addresses.blurbDB, dialOpts)
	if err != nil {
		log.Printf("SERVER: Failed to connect to BlurbDB at %s", configuration.addresses.blurbDB)
	}
	configuration.clients.blurbDB = blurb.NewBlurbDBClient(configuration.connections.blurb)
	log.Printf("SERVER: Successfully created a connection to blurb db")

	// User database
	configuration.connections.user, err = grpc.Dial(configuration.addresses.userDB, dialOpts)
	if err != nil {
		log.Printf("SERVER: Failed to connect to UserDB at %s", configuration.addresses.userDB)
	}
	configuration.clients.userDB = user.NewUserDBClient(configuration.connections.user)
	log.Printf("SERVER: Successfully created a connection to user db")

	// Subscription databases
	configuration.clients.subDBs = make(map[string]subscription.SubscriptionDBClient, 3)
	configuration.connections.subs = make(map[string]*grpc.ClientConn, 3)

	for idx := 1; idx < 4; idx++ {
		addr := getAddress("SUB" + strconv.Itoa(idx))
		configuration.connections.subs[addr], err = grpc.DialContext(context.Background(), addr, dialOpts)
		if err != nil {
			log.Printf("SERVER: Failed to connect to SubDB at %s", addr)
		}
		configuration.clients.subDBs[addr] = subscription.NewSubscriptionDBClient(configuration.connections.subs[addr])
	}
	log.Printf("SERVER: Successfully created a connection to sub dbs")
}

func initializeConnectors() {
	// init the mutexes
	configuration.muts.blurbMut = &sync.Mutex{}
	configuration.muts.subMut = &sync.Mutex{}
	configuration.muts.userMut = &sync.Mutex{}

	// gather the addresses for the gRPC servers
	gatherAddresses()

	// dial the servers
	dialServers()

	// start maintaining a history of who the leader is
	go configuration.maintainConnections()

	log.Printf("SERVER: Initialized! All gRPC connections secured\n")
}
