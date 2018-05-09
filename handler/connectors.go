package handler

import (
	"context"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc/connectivity"

	"github.com/adamsanghera/blurber-protobufs/dist/blurb"
	"github.com/adamsanghera/blurber-protobufs/dist/common"
	"github.com/adamsanghera/blurber-protobufs/dist/subscription"
	"github.com/adamsanghera/blurber-protobufs/dist/user"
	"google.golang.org/grpc"
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

type vote struct {
	err  error
	info *common.ServerInfo
}

// updateLeaders performs a parallelized poll of the gRPC servers, to determine
// which is the leader
func (c *config) updateLeaders() {
	c.muts.subMut.Lock()
	defer c.muts.subMut.Unlock()
	voteChannel := make(chan vote)
	for addr, client := range c.clients.subDBs {

		// Spawn a function to get a vote from every client
		go func(vChan chan vote, addr string, subclient subscription.SubscriptionDBClient) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
			info, err := subclient.Leader(ctx, &common.Empty{})
			if info == nil {
				vChan <- vote{err, &common.ServerInfo{Address: addr[len(addr)-5:len(addr)-1] + "1", Index: -1}}
			} else {
				vChan <- vote{err, info}
			}
			cancel()
		}(voteChannel, addr, client)

	}

	voteBox := make(map[string]int)

	for range c.clients.subDBs {
		vote := <-voteChannel
		if vote.err != nil {
			voterAddr := vote.info.Address

			log.Printf("SERVER CONN: Failed to contact %s, with msg {%s}", voterAddr, vote.err.Error())
			// Only call for a view change if the absent voter is the leader
			// log.Printf("Testing: %s == %s", voterAddr[len(voterAddr)-4:len(voterAddr)-1], c.leaders.sub[1:4])
			if c.leaders.sub != "" && voterAddr[len(voterAddr)-5:len(voterAddr)-1] == c.leaders.sub[:4] {
				c.leaders.sub = ""
				// prompt view change
				for addr, client := range c.clients.subDBs {
					if c.connections.subs[addr].GetState() == connectivity.Ready {
						log.Printf("Prompting view change at %s", addr)
						_, err := client.PromptViewChange(context.Background(), &common.Empty{})
						if err != nil {
							log.Printf("SERVER CONN: Failed to prompt view change at %s, err msg {%s}", addr, err.Error())
						} else {
							return
						}
					}
				}
			}
		} else {
			if _, exists := voteBox[vote.info.Address]; !exists {
				voteBox[vote.info.Address] = 0
			}
			voteBox[vote.info.Address]++
		}
	}

	winner := ""
	winningTally := 0
	for key, value := range voteBox {
		if value > winningTally {
			winner = key
		}
	}

	log.Printf("Vote results: %v", voteBox)

	c.leaders.sub = winner
}

// maintainConnections calls updateLeaders periodically
func (c *config) maintainConnections() {
	epoch := 0
	for {
		c.updateLeaders()
		c.muts.subMut.Lock()
		if c.leaders.sub == "" {
			c.muts.subMut.Unlock()
			c.updateLeaders()
		} else {
			c.muts.subMut.Unlock()
		}
		log.Printf("SERVER: Leader for epoch %d is %s", epoch, c.leaders.sub)
		epoch++
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

	addr := c.addresses.blurbDB[:len(c.addresses.blurbDB)-5] + c.leaders.sub[:len(c.leaders.sub)-1] + "0"

	if c.connections.subs[addr].GetState() != connectivity.Ready {
		log.Printf("SERVER: Leader connection %s is not ready... calling an election", c.leaders.sub)

		// Need to release lock for the election
		c.muts.subMut.Unlock()
		c.updateLeaders()
		c.muts.subMut.Lock()
	}

	return c.clients.subDBs[addr]
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
	// dialOpts := grpc.WithKeepaliveParams(keepalive.ClientParameters{
	// 	Time:                5 * time.Second,
	// 	Timeout:             2 * time.Second,
	// 	PermitWithoutStream: false,
	// })
	dialOpts := grpc.WithInsecure()

	// Blurb database
	var err error
	configuration.connections.blurb, err = grpc.Dial(configuration.addresses.blurbDB, dialOpts)
	if err != nil {
		log.Printf("SERVER: Failed to connect to BlurbDB at %s", configuration.addresses.blurbDB)
	} else {
		configuration.clients.blurbDB = blurb.NewBlurbDBClient(configuration.connections.blurb)
		log.Printf("SERVER: Successfully created a connection to BlurbDB at %s", configuration.addresses.blurbDB)
	}

	// User database
	configuration.connections.user, err = grpc.Dial(configuration.addresses.userDB, dialOpts)
	if err != nil {
		log.Printf("SERVER: Failed to connect to UserDB at %s", configuration.addresses.userDB)
	} else {
		configuration.clients.userDB = user.NewUserDBClient(configuration.connections.user)
		log.Printf("SERVER: Successfully created a connection to UserDB at %s", configuration.addresses.userDB)
	}

	// Subscription databases
	configuration.clients.subDBs = make(map[string]subscription.SubscriptionDBClient, 3)
	configuration.connections.subs = make(map[string]*grpc.ClientConn, 3)

	for idx := 1; idx < 4; idx++ {
		addr := getAddress("SUB" + strconv.Itoa(idx))
		configuration.connections.subs[addr], err = grpc.DialContext(context.Background(), addr, dialOpts)
		if err != nil {
			log.Printf("SERVER: Failed to connect to SubDB at %s", addr)
		} else {
			configuration.clients.subDBs[addr] = subscription.NewSubscriptionDBClient(configuration.connections.subs[addr])
			log.Printf("SERVER: Successfully created a connection to SubDB %s", addr)
		}
	}
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

	// start maintaining a history of who the leader is, sending messages
	go configuration.maintainConnections()

	time.Sleep(1 * time.Second)

	log.Printf("SERVER: Finished Initialization! All gRPC connections secured:\n %v", configuration)
}
