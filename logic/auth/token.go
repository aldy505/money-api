package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var ErrTokenParse = errors.New("error parsing jwt token")
var ErrTokenInvalid = errors.New("invalid jwt token")
var ErrTokenNotAcceptable = errors.New("jwt token not acceptable")

func GenerateJWT(secret []byte, user User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["nbf"] = time.Now().Unix()
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	claims["email"] = user.Email
	claims["id"] = user.ID
	claims["iss"] = os.Getenv("JWT_ISSUER")
	claims["sub"] = "PYPL_TKN"

	s, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return s, nil
}

func VerifyJWT(secret []byte, userToken string) (User, error) {
	token, err := jwt.Parse(userToken, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrTokenParse
		}
		return secret, nil
	})
	if err != nil {
		return User{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return User{}, ErrTokenInvalid
	}

	if claims["iss"] != os.Getenv("JWT_ISSUER") {
		return User{}, ErrTokenNotAcceptable
	}

	var user User
	user.Email = claims["email"].(string)
	id, _ := claims["id"].(float64)
	user.ID = int(id)
	return user, nil
}
