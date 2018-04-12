package registration

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
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
	userSet  map[int32]user
	userID   map[string]int32
	pwdMap   map[int32]string
	tokenMap map[int32]Token

	uidCounter int32
}

// NewLocalLedger initializes and returns a new LocalLedger instance.
func NewLocalLedger() *LocalLedger {
	log.Printf("AUTH-LEDGER: Initializing")
	return &LocalLedger{
		userSet:    make(map[int32]user),
		userID:     make(map[string]int32),
		pwdMap:     make(map[int32]string),
		tokenMap:   make(map[int32]Token),
		uidCounter: 0,
	}
}

// AddNewUser integrates a new user into the ledger, with 3 operations.
//   (1) Assignment of a permanent UID
//   (2) Instantiation of a new User object, associated with the UID
//   (3) Storage of a provided password, associated with the UID.
// The ledger's UID counter is then (TODO: thread-safely) incremented.
func (lul *LocalLedger) AddNewUser(name string, pwd string) error {
	log.Printf("AUTH-LEDGER: Registering new user %s", name)

	if _, exists := lul.userID[name]; exists {
		return errors.New("A user with that name already exists")
	}

	lul.userSet[lul.uidCounter] = user{
		Name: name,
		UID:  lul.uidCounter,
	}
	lul.userID[name] = lul.uidCounter
	lul.pwdMap[lul.uidCounter] = pwd

	// TODO: Make threadsafe
	lul.uidCounter++

	return nil
}

// LogIn performs a username-check, followed by a password-check.
// If both checks pass, then a new token is allocated to associated UID, and returned.
func (lul *LocalLedger) LogIn(uname string, pwd string) (string, error) {
	log.Printf("AUTH-LEDGER: Logging-in %s", uname)

	// username-check
	id, err := lul.GetUserID(uname)
	if err != nil {
		return "", err
	}

	// password-check
	if lul.pwdMap[id] != pwd {
		return "", errors.New("Bad password for " + uname)
	}

	return lul.allocateNewToken(id), nil
}

func (lul *LocalLedger) allocateNewToken(uid int32) string {
	// Generate a random token of fixed-length
	bitString := make([]byte, 256)
	_, err := rand.Read(bitString)
	if err != nil {
		panic(err)
	}
	token := hex.EncodeToString(bitString)

	// Associate the token with the given uid
	lul.tokenMap[uid] = Token{token: token, creationDate: time.Now()}

	return token
}

// CheckIn performs a username-check and a token-check.
// If both checks pass, then a new token is validated and returned
func (lul *LocalLedger) CheckIn(uname string, token string) (string, error) {
	log.Printf("AUTH-LEDGER: Checking-in %s", uname)

	// username-check
	id, err := lul.GetUserID(uname)
	if err != nil {
		return "", err
	}

	// token-check: part 1
	if lul.tokenMap[id].token != token {
		return "", errors.New("Bad token for " + uname)
	}

	// token-check: part 2
	if time.Since(lul.tokenMap[id].creationDate) > 30*time.Minute {
		return "", errors.New("Session expried for " + uname)
	}

	return lul.allocateNewToken(id), nil
}

func (lul *LocalLedger) CheckOut(uname string, token string) error {
	log.Printf("AUTH-LEDGER: Checking-out %s", uname)

	// username-check
	id, err := lul.GetUserID(uname)
	if err != nil {
		return err
	}

	// token-check
	if lul.tokenMap[id].token != token {
		return errors.New("Bad token for " + uname)
	}

	delete(lul.tokenMap, id)

	return nil
}

// GetUserID retrieves the permanent UID associated with uname.
func (lul *LocalLedger) GetUserID(uname string) (int32, error) {
	log.Printf("AUTH-LEDGER: Retrieving UID for %s", uname)

	id, ok := lul.userID[uname]
	if !ok {
		return -1, errors.New("No record of " + uname + " exists")
	}
	return id, nil
}

// Remove eliminates all history of uname, including its UID.
func (lul *LocalLedger) Remove(uname string) error {
	log.Printf("AUTH-LEDGER: Removing all history of %s", uname)

	// username-check
	id, err := lul.GetUserID(uname)
	if err != nil {
		return err
	}

	delete(lul.userID, uname)
	delete(lul.userSet, id)
	delete(lul.pwdMap, id)
	delete(lul.tokenMap, id)
	return nil
}
