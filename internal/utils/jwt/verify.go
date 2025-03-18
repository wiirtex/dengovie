package jwt

import (
	"dengovie/internal/domain"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"time"
)

type JWTData struct {
	UserID string `json:"user_id"`
}

var (
	errInvalidToken = errors.New("invalid token")
	errExpiredToken = errors.New("expired token")
)

func VerifyJWT(tok string) (JWTData, error) {

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
		return JWTData{}, errInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {

		exp, err := claims.GetExpirationTime()
		if err != nil {
			log.Println("claims.GetExpirationTime error:", err)
			return JWTData{}, errInvalidToken
		}

		notBefore, err := claims.GetNotBefore()
		if err != nil {
			log.Println("claims.GetNotBefore error:", err)
			return JWTData{}, errInvalidToken
		}

		if notBefore.After(time.Now()) || exp.Before(time.Now()) {
			log.Println("expired token:", err)
			return JWTData{}, errExpiredToken
		}

		fmt.Printf("claims: %v\n", claims[domain.UserIDKey])
		userID, ok := claims[domain.UserIDKey].(string)
		if !ok {
			log.Println("userID is not a string")
			return JWTData{}, errInvalidToken
		}
		return JWTData{
			UserID: userID,
		}, nil
	}

	return JWTData{}, fmt.Errorf("token.Claims are not jwt.MapClaims")
}
