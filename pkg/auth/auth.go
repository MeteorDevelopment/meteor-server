package auth

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"golang.org/x/exp/slices"
	"strings"
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

var jwtKey []byte

var tokenCount = 0
var tokens = make(map[ksuid.KSUID][]int)
var mu = sync.RWMutex{}

var confirmEmails = make(map[ksuid.KSUID]ConfirmEmailStruct)
var cetMu = sync.RWMutex{}

func Init() {
	jwtKey = make([]byte, 36)
	rand.Read(jwtKey)
}

func Register(username string, email string, password string) error {
	if username == "" || strings.ContainsRune(username, ' ') || email == "" || strings.ContainsRune(email, ' ') {
		return errors.New("Invalid username or email.")
	}

	if err := ValidatePassword(password); err != nil {
		return err
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

	core.SendEmail(email, "Confirm email to register", "To complete the registration go to https://meteorclient.com/confirm?token="+token.String()+" . The link is valid for 15 minutes.")
	return nil
}

// ValidatePassword checks the length in bytes, bcrypt can't hash passwords longer than 72 bytes
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("Password must be at least 8 characters long.")
	}

	if len(password) > 72 {
		return errors.New("Password can't be longer than 72 characters.")
	}

	return nil
}

func ConfirmEmail(token ksuid.KSUID) bool {
	cetMu.Lock()
	clearConfirmEmails()
	confirmEmail, exists := confirmEmails[token]
	if !exists {
		cetMu.Unlock()
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
		mu.Unlock()
		return "", err
	}

	token, err := jose.Sign(string(bytes), jose.HS256, jwtKey)
	if err != nil {
		mu.Unlock()
		return "", err
	}

	tokens[account.ID] = append(tokens[account.ID], tokenCount)
	tokenCount++

	mu.Unlock()
	return token, nil
}

func Logout(token string, id ksuid.KSUID) error {
	bytes, _, err := jose.Decode(token, jwtKey)
	if err != nil {
		return err
	}

	var claims Claims
	err = json.Unmarshal([]byte(bytes), &claims)
	if err != nil {
		return err
	}

	mu.Lock()
	tokenIds, exists := tokens[id]
	if exists {
		for i := 0; i < len(tokenIds); i++ {
			if tokenIds[i] == claims.TokenID {
				tokens[id] = append(tokenIds[:i], tokenIds[i+1:]...)
				break
			}
		}
	}
	mu.Unlock()

	return nil
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
	validTokenIds, exists := tokens[claims.AccountID]
	mu.RUnlock()

	if exists && slices.Contains(validTokenIds, claims.TokenID) {
		return claims.AccountID, nil
	}

	return ksuid.Nil, errors.New("invalid token")
}

func Invalidate(id ksuid.KSUID) {
	delete(tokens, id)
}
