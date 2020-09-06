package common

import (
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

const (
	tokenValidity = time.Hour * 4
)

var SigningJWTError = errors.New("error signing jtw token with key")
var KeyFuncJWTError = errors.New("error with KeyFunc")
var ParsingJWTError = errors.New("error with JWT parsing")
var InvalidJWTTokenError = errors.New("error the token is invalid")

func GenerateJWT(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["sub"] = username
	claims["nbf"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(tokenValidity).Unix()
	claims["iss"] = "auth microservice"

	key := os.Getenv("JWT_KEY")
	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		return "", SigningJWTError
	}
	return signedToken, nil
}

func VerifyJWT(token string) (string, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", KeyFuncJWTError
		}
		key := os.Getenv("JWT_KEY")
		return []byte(key), nil
	})
	if err != nil {
		return "", ParsingJWTError
	}

	if parsedToken.Valid {
		claims := parsedToken.Claims.(jwt.MapClaims)
		return claims["sub"].(string), nil
	}
	return "", InvalidJWTTokenError
}
