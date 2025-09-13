package jwt

import (
	"dengovie/internal/domain"
	"dengovie/internal/web"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	errInvalidToken = errors.New("invalid token")
	errExpiredToken = errors.New("expired token")
)

func (j *jwtProcessor) VerifyJWT(tok string) (map[web.JWTKey]any, error) {
	initOnce()

	token, err := jwt.Parse(tok, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			log.Printf("Unexpected signing method: %v\n", token.Header["alg"])
			return nil, errInvalidToken
		}

		return jwtKey, nil
	})
	if err != nil {
		log.Println("jwt.Parse:", err)
		return nil, errInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {

		exp, err := claims.GetExpirationTime()
		if err != nil {
			log.Println("claims.GetExpirationTime error:", err)
			return nil, errInvalidToken
		}

		notBefore, err := claims.GetNotBefore()
		if err != nil {
			log.Println("claims.GetNotBefore error:", err)
			return nil, errInvalidToken
		}

		if notBefore.After(time.Now()) || exp.Before(time.Now()) {
			log.Println("expired token:", err)
			return nil, errExpiredToken
		}

		userID, ok := claims[domain.UserIDKey].(float64)
		if !ok {
			log.Println("userID is not float64")
			return nil, errInvalidToken
		}

		return map[web.JWTKey]any{
			web.JWTUserIDKey: domain.UserID(userID),
		}, nil
	}

	return nil, fmt.Errorf("token.Claims are not jwt.MapClaims")
}
