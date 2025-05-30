package aaa

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "phoenix"
const adminRole = "superuser"

type Logger interface {
	Error(msg string, keysAndValues ...interface{})
}

type AAA struct {
	users    map[string]string
	tokenTTL time.Duration
	log      Logger
}

func New(tokenTTL time.Duration, log Logger) (AAA, error) {
	const adminUser = "ADMIN_USER"
	const adminPass = "ADMIN_PASSWORD"

	user, exists := os.LookupEnv(adminUser)
	if !exists {
		return AAA{}, fmt.Errorf("could not get admin user from enviroment")
	}
	password, exists := os.LookupEnv(adminPass)
	if !exists {
		return AAA{}, fmt.Errorf("could not get admin password from enviroment")
	}

	return AAA{
		users:    map[string]string{user: password},
		tokenTTL: tokenTTL,
		log:      log,
	}, nil
}

func (a AAA) Login(name, password, sub string) (string, error) {
	if name == "" {
		return "", errors.New("empty user")
	}
	savedPass, ok := a.users[name]
	if !ok {
		return "", errors.New("unknown user")
	}
	if savedPass != password {
		return "", errors.New("wrong password")
	}

	if sub == "" {
		sub = adminRole
	}

	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  sub,
		"name": name,
		"exp":  jwt.NewNumericDate(time.Now().Add(a.tokenTTL)),
	})

	return token.SignedString([]byte(secretKey))
}

func (a AAA) Verify(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		a.log.Error("cannot parse token", "error", err)
		return fmt.Errorf("cannot parse token")
	}
	if !token.Valid {
		a.log.Error("token is invalid")
		return errors.New("token is invalid")
	}
	subject, err := token.Claims.GetSubject()
	if err != nil {
		a.log.Error("no subject", "error", err)
		return errors.New("incomplete token")
	}
	if subject != adminRole {
		a.log.Error("not admin", "subject", subject)
		return errors.New("not authorized")
	}
	return nil
}
