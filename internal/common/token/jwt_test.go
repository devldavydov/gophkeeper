package token

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWTToken(t *testing.T) {
	secret := []byte("foobar")
	userID := int64(1)

	for _, tt := range []struct {
		name  string
		exp   time.Duration
		isErr bool
	}{
		{name: "valid token", exp: 100 * time.Second},
		{name: "expired token", exp: 1 * time.Second, isErr: true},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			token, err := NewJWTForUser(userID, tt.exp, secret)
			assert.NoError(t, err)
			assert.NotEqual(t, "", token)

			time.Sleep(2 * time.Second)

			tokenUserID, err := GetUserFromJWT(token, secret)
			if tt.isErr {
				assert.ErrorIs(t, err, ErrInvalidToken)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, strconv.FormatInt(userID, 10), tokenUserID)
			}
		})
	}
}
