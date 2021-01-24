package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var loc, _ = time.LoadLocation("America/Sao_Paulo")
var jwtKey = "1.2.3.4.5"

// Generate a token jwt, user key = "" to assume the default value
func Generate(key string, assingWithDefaultKey bool, id string, origin string) (string, error) {
	if assingWithDefaultKey {
		key = jwtKey
	}

	var (
		token  = jwt.New(jwt.SigningMethodHS256)
		claims = token.Claims.(jwt.MapClaims)
		now    = time.Now().In(loc)
	)

	claims["iss"] = origin                             // issue - origin
	claims["iat"] = now                                // issueAt - timestamp de quando foi geral
	claims["exp"] = now.Add(time.Hour * 24 * 7).Unix() // expiration - timestamp de quando deve expirar
	claims["sub"] = id                                 // subject - geralmente id do usuário

	hash, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return hash, nil
}

// Verify if a hash is a validate jwt, use key = "" to assume the default value
func Verify(key string, assingWithDefaultKey bool, hash string, recoveryClaims bool) (r bool, c map[string]interface{}, e error) {
	if assingWithDefaultKey {
		key = jwtKey
	}

	r = false
	e = errors.New("Token inválido")

	if hash == "" {
		return
	}

	token, err := jwt.Parse(hash, func(hash *jwt.Token) (interface{}, error) {
		if _, ok := hash.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Erro ao validar token")
		}
		return []byte(key), nil
	})

	if err != nil {
		e = err
		return
	}

	now := time.Now().In(loc).Unix()
	exp := token.Claims.(jwt.MapClaims)["exp"].(float64)

	if int64(exp) > now && token.Valid {
		r, e = true, nil
	}

	if recoveryClaims {
		c = token.Claims.(jwt.MapClaims)
	}

	return
}
