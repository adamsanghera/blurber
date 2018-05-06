package registration

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"sync"
	"time"
)

// Ledger is the interfance for a service that keeps track
// of user passwords, tokens, screennames, and permanent IDs.
type Ledger interface {
	Add(u user, pwd string) error
	Remove(uname string, pwd string) error
	LogIn(uname string, pwd string) (error, string) // returns a token
	LogOut(uname string, pwd string) error
	CheckIn(uname string, token string) (error, string) // returns a token
	CheckOut(uname string, token string) error
}

type Token struct {
	creationDate time.Time
	token        string
}

// LocalLedger is an implementation of Ledger, which uses in-memory data
// structures to support Ledger's functionality.
type LocalLedger struct {
	userSet  sync.Map
	userID   sync.Map
	pwdMap   sync.Map
	tokenMap sync.Map

	uidCounter int32
	uidMut     sync.Mutex
}

// NewLocalLedger initializes and returns a new LocalLedger instance.
func NewLocalLedger() *LocalLedger {
	log.Printf("AUTH-LEDGER: Initializing")
	return &LocalLedger{
		userSet:    sync.Map{},
		userID:     sync.Map{},
		pwdMap:     sync.Map{},
		tokenMap:   sync.Map{},
		uidCounter: 0,
	}
}

func (ll *LocalLedger) measureUID() int32 {
	ll.uidMut.Lock()
	ret := ll.uidCounter
	ll.uidCounter++
	ll.uidMut.Unlock()
	return ret
}

// AddNewUser integrates a new user into the ledger, with 3 operations.
//   (1) Assignment of a permanent UID
//   (2) Instantiation of a new User object, associated with the UID
//   (3) Storage of a provided password, associated with the UID.
// The ledger's UID counter is then (TODO: thread-safely) incremented.
func (ll *LocalLedger) AddNewUser(name string, pwd string) error {
	log.Printf("AUTH-LEDGER: Registering new user %s", name)

	if _, exists := ll.userID.Load(name); exists {
		return errors.New("A user with that name already exists")
	}

	id := ll.measureUID()

	ll.userSet.Store(id, user{
		Name: name,
		UID:  id,
	})

	ll.userID.Store(name, id)
	ll.pwdMap.Store(id, pwd)

	return nil
}

// LogIn performs a username-check, followed by a password-check.
// If both checks pass, then a new token is allocated to associated UID, and returned.
func (ll *LocalLedger) LogIn(uname string, pwd string) (string, error) {
	log.Printf("AUTH-LEDGER: Logging-in %s", uname)

	// username-check
	id, err := ll.GetUserID(uname)
	if err != nil {
		return "", err
	}

	// password-check
	pass, exists := ll.pwdMap.Load(id)
	if !exists {
		panic("Password for user doesn't exist!")
	}

	if pass != pwd {
		return "", errors.New("Bad password for " + uname)
	}

	return ll.allocateNewToken(id), nil
}

func (ll *LocalLedger) allocateNewToken(uid int32) string {
	// Generate a random token of fixed-length
	bitString := make([]byte, 256)
	_, err := rand.Read(bitString)
	if err != nil {
		panic(err)
	}
	token := hex.EncodeToString(bitString)

	// Associate the token with the given uid
	ll.tokenMap.Store(uid, Token{token: token, creationDate: time.Now()})

	return token
}

func (ll *LocalLedger) getToken(uid int32) Token {
	tok, exists := ll.tokenMap.Load(uid)
	if !exists {
		return Token{}
	}

	token, ok := tok.(Token)
	if !ok {
		panic("Bad token in DB")
	}

	return token
}

// CheckIn performs a username-check and a token-check.
// If both checks pass, then a new token is validated and returned
func (ll *LocalLedger) CheckIn(uname string, token string) (string, error) {
	log.Printf("AUTH-LEDGER: Checking-in %s", uname)

	// username-check
	id, err := ll.GetUserID(uname)
	if err != nil {
		return "", err
	}

	internalToken := ll.getToken(id)

	// token-check: part 1
	if internalToken.token != token {
		return "", errors.New("Bad token for " + uname)
	}

	// token-check: part 2
	if time.Since(internalToken.creationDate) > 30*time.Minute {
		return "", errors.New("Session expried for " + uname)
	}

	return ll.allocateNewToken(id), nil
}

func (ll *LocalLedger) CheckOut(uname string, token string) error {
	log.Printf("AUTH-LEDGER: Checking-out %s", uname)

	// username-check
	id, err := ll.GetUserID(uname)
	if err != nil {
		return err
	}

	internalToken := ll.getToken(id)

	// token-check
	if internalToken.token != token {
		return errors.New("Bad token for " + uname)
	}

	ll.tokenMap.Delete(id)

	return nil
}

// GetUserID retrieves the permanent UID associated with uname.
func (ll *LocalLedger) GetUserID(uname string) (int32, error) {
	log.Printf("AUTH-LEDGER: Retrieving UID for %s", uname)

	id, ok := ll.userID.Load(uname)
	if !ok {
		return -1, errors.New("No record of " + uname + " exists")
	}

	intID, ok := id.(int32)
	if !ok {
		panic("Bad id value stored in map")
	}

	return intID, nil
}

// Remove eliminates all history of uname, including its UID.
func (ll *LocalLedger) Remove(uname string) error {
	log.Printf("AUTH-LEDGER: Removing all history of %s", uname)

	// username-check
	id, err := ll.GetUserID(uname)
	if err != nil {
		return err
	}

	ll.userSet.Delete(uname)
	ll.userID.Delete(id)
	ll.pwdMap.Delete(id)
	ll.tokenMap.Delete(id)
	return nil
}
