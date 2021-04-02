package auth

import (
	"encoding/json"
	"errors"
	"sync"

	jose "github.com/dvsekhvalnov/jose2go"
	"meteor-server/pkg/db"
)

type Claims struct {
	TokenID   int
	AccountID string
}

// TODO: Randomly generate the key on each startup
var jwtKey = []byte{97,48,97,50,97,98,100,56,45,54,49,54,50,45,52,49,99,51,45,56,51,100,54,45,49,99,102,53,53,57,98,52,54,97,102,99}

var tokenCount = 0
var tokens = make(map[string]int)
var mu = sync.RWMutex{}

func Login(name string, password string) (string, error) {
	if name == "" || password == "" {
		return "", errors.New("wrong name or password")
	}

	account, err := db.GetAccountWithUsernameOrEmail(name)
	if err != nil {
		return "", errors.New("wrong name or password")
	}

	if password != account.Password {
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

func Logout(accountId string) {
	mu.Lock()
	delete(tokens, accountId)
	mu.Unlock()
}

func IsTokenValid(token string) (string, error) {
	bytes, _, err := jose.Decode(token, jwtKey)
	if err != nil {
		return "", err
	}

	var claims Claims
	err = json.Unmarshal([]byte(bytes), &claims)
	if err != nil {
		return "", err
	}

	mu.RLock()
	validTokenId, exists := tokens[claims.AccountID]
	mu.RUnlock()

	if exists && claims.TokenID == validTokenId {
		return claims.AccountID, nil
	}

	return "", errors.New("invalid token")
}
