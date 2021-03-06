package auth

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	jose "github.com/dvsekhvalnov/jose2go"
	"github.com/segmentio/ksuid"
	"meteor-server/pkg/core"
	"meteor-server/pkg/db"
)

type Claims struct {
	TokenID   int
	AccountID ksuid.KSUID
}

type ConfirmEmailStruct struct {
	Token    ksuid.KSUID
	Username string
	Email    string
	Password string
	Time     time.Time
}

// TODO: Randomly generate the key on each startup
var jwtKey = []byte{97, 48, 97, 50, 97, 98, 100, 56, 45, 54, 49, 54, 50, 45, 52, 49, 99, 51, 45, 56, 51, 100, 54, 45, 49, 99, 102, 53, 53, 57, 98, 52, 54, 97, 102, 99}

var tokenCount = 0
var tokens = make(map[ksuid.KSUID]int)
var mu = sync.RWMutex{}

var confirmEmails = make(map[ksuid.KSUID]ConfirmEmailStruct)
var cetMu = sync.RWMutex{}

func Register(username string, email string, password string) error {
	if username == "" || email == "" || password == "" {
		return errors.New("Empty username, email or password.")
	}

	if !core.IsEmailValid(email) {
		return errors.New("Invalid email.")
	}

	_, err := db.GetAccountWithUsername(username)
	if err == nil {
		return errors.New("Account with this username already exists.")
	}

	_, err = db.GetAccountWithEmail(email)
	if err == nil {
		return errors.New("Account with this email already exists.")
	}

	token := ksuid.New()

	cetMu.Lock()
	clearConfirmEmails()
	confirmEmails[token] = ConfirmEmailStruct{Token: token, Username: username, Email: email, Password: password, Time: time.Now()}
	cetMu.Unlock()

	core.SendEmail(email, "Confirm email to register", "To complete the registration go to https://meteorclient.com/confirm?token=" + token.String() + ". The link is valid for 15 minutes.")
	return nil
}

func ConfirmEmail(token ksuid.KSUID) bool {
	cetMu.Lock()
	clearConfirmEmails()
	confirmEmail, exists := confirmEmails[token]
	if !exists {
		return false
	}

	delete(confirmEmails, token)
	cetMu.Unlock()

	err := db.NewAccount(confirmEmail.Username, confirmEmail.Email, confirmEmail.Password)
	return err == nil
}

func clearConfirmEmails() {
	now := time.Now()

	for token, confirmEmail := range confirmEmails {
		if now.Sub(confirmEmail.Time).Minutes() > 15 {
			delete(confirmEmails, token)
		}
	}
}

func Login(name string, password string) (string, error) {
	if name == "" || password == "" {
		return "", errors.New("wrong name or password")
	}

	account, err := db.GetAccountWithUsernameOrEmail(name)
	if err != nil {
		return "", errors.New("wrong name or password")
	}

	if !account.PasswordMatches(password) {
		return "", errors.New("wrong name or password")
	}

	mu.Lock()

	bytes, err := json.Marshal(Claims{TokenID: tokenCount, AccountID: account.ID})
	if err != nil {
		return "", err
	}

	token, err := jose.Sign(string(bytes), jose.HS256, jwtKey)
	if err != nil {
		return "", err
	}

	tokens[account.ID] = tokenCount
	tokenCount++

	mu.Unlock()
	return token, nil
}

func Logout(id ksuid.KSUID) {
	mu.Lock()
	delete(tokens, id)
	mu.Unlock()
}

func IsTokenValid(token string) (ksuid.KSUID, error) {
	bytes, _, err := jose.Decode(token, jwtKey)
	if err != nil {
		return ksuid.Nil, err
	}

	var claims Claims
	err = json.Unmarshal([]byte(bytes), &claims)
	if err != nil {
		return ksuid.Nil, err
	}

	mu.RLock()
	validTokenId, exists := tokens[claims.AccountID]
	mu.RUnlock()

	if exists && claims.TokenID == validTokenId {
		return claims.AccountID, nil
	}

	return ksuid.Nil, errors.New("invalid token")
}
