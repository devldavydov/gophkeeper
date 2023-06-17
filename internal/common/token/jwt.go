// Package token contains utils to work with user tokens.
package token

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid token")

// NewJWTForUser creates new JWT token for user.
func NewJWTForUser(userID int64, expirationInterval time.Duration, secret []byte) (string, error) {
	claims := &jwt.RegisteredClaims{
		ID:        strconv.FormatInt(userID, 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expirationInterval)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func GetUserFromJWT(tokenString string, secret []byte) (int64, error) {
	var token *jwt.Token
	var err error

	token, err = jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}

			return secret, nil
		},
		jwt.WithTimeFunc(func() time.Time { return time.Now().UTC() }),
	)

	if err != nil {
		return 0, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		var userID int64

		userID, err = strconv.ParseInt(claims.ID, 10, 64)
		if err != nil {
			return 0, ErrInvalidToken
		}
		return userID, nil
	}

	return 0, ErrInvalidToken
}
