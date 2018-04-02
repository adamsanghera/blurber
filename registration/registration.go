package registration

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"time"
)

type UserLedger interface {
	Add(u User, pwd string) error
	Remove(uname string, pwd string) error
	LogIn(uname string, pwd string) (error, string) // returns a token
	LogOut(uname string, pwd string) error
	Authorize(uname string, token string) (error, string) // returns a token
}

type Token struct {
	creationDate time.Time
	token        string
}

type LocalUserLedger struct {
	userSet  map[string]User
	userID   map[string]string
	pwdMap   map[string]string
	tokenMap map[string]Token
}

func NewLocalLedger() *LocalUserLedger {
	return &LocalUserLedger{
		userSet:  make(map[string]User),
		userID:   make(map[string]string),
		pwdMap:   make(map[string]string),
		tokenMap: make(map[string]Token),
	}
}

// Not threadsafe
func (lul *LocalUserLedger) Add(u User, pwd string) error {
	log.Printf("REGISTRATION: Adding user %s with pwd %s", u.Name, pwd)
	lul.userSet[u.UID] = u
	lul.userID[u.Name] = u.UID
	lul.pwdMap[u.UID] = pwd

	return nil
}

// Not threadsafe
func (lul *LocalUserLedger) LogIn(uname string, pwd string) (error, string) {
	log.Printf("LOGIN: Attempting user %s with pwd %s", uname, pwd)

	id, ok := lul.userID[uname]
	if !ok {
		return errors.New("No record of " + uname + " exists"), ""
	}

	if lul.pwdMap[id] != pwd {
		return errors.New("Bad password for " + uname), ""
	}

	return nil, lul.upsertToken(id)
}

func (lul *LocalUserLedger) upsertToken(id string) string {
	bitString := make([]byte, 256)
	_, err := rand.Read(bitString)
	if err != nil {
		panic(err)
	}

	token := hex.EncodeToString(bitString)

	lul.tokenMap[id] = Token{token: token, creationDate: time.Now()}

	return token
}

func (lul *LocalUserLedger) Authorize(uname string, token string) (error, string) {
	id, ok := lul.userID[uname]
	if !ok {
		return errors.New("No record of " + uname + " exists"), ""
	}

	if lul.tokenMap[id].token != token {
		return errors.New("Bad token for " + uname), ""
	}

	if time.Since(lul.tokenMap[id].creationDate) > 30*time.Second {
		return errors.New("Session expried for " + uname), ""
	}

	return nil, lul.upsertToken(id)
}

func (lul *LocalUserLedger) GetUsrID(uname string) (error, string) {
	log.Printf("LOGIN: Getting user %s", uname)

	id, ok := lul.userID[uname]
	if !ok {
		return errors.New("No record of " + uname + " exists"), ""
	}
	return nil, id
}
