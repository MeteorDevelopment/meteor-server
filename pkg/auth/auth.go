package auth

import (
	"errors"
	"sync"

	"github.com/dgrijalva/jwt-go/v4"
	"meteor-server/pkg/db"
)

type Claims struct {
	TokenID   int
	AccountID string

	jwt.StandardClaims
}

// TODO: Randomly generate the key on each startup
var jwtKey = []byte("OMGF a key !!!?!?!?")

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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{TokenID: tokenCount, AccountID: account.ID})
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	tokens[account.ID] = tokenCount
	tokenCount++

	mu.Unlock()
	return tokenStr, nil
}

func Logout(accountId string) {
	mu.Lock()
	delete(tokens, accountId)
	mu.Unlock()
}

func IsTokenValid(tokenStr string) (string, error) {
	token, _ := jwt.ParseWithClaims(tokenStr, Claims{}, func(token *jwt.Token) (interface{}, error) {
		return token, nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		mu.RLock()
		validTokenId, exists := tokens[claims.AccountID]
		mu.RUnlock()

		if exists && claims.TokenID == validTokenId {
			return claims.AccountID, nil
		}
	}

	return "", errors.New("invalid token")
}
