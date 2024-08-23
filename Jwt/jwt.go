package Jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var refreshSecretKey = []byte("refresh-secret-key")
var accessSecretKey = []byte("access-secret-key")

func CreateRefreshToken(username string) (string, error) {

	refreshTokenClaim := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
		})

	refreshToken, err := refreshTokenClaim.SignedString(refreshSecretKey)
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func CreateAccessToken(refreshT string) (string, error) {

	accessTokenClaim := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": refreshT,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	accessToken, err := accessTokenClaim.SignedString(accessSecretKey)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func AccessTokenValidity(tokenString string) (*jwt.Token, error) {

	if tokenString == "" {
		return nil, fmt.Errorf("missing token in header")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return accessSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func RefreshTokenValidity(tokenString string) (*jwt.Token, error) {

	if tokenString == "" {
		return nil, fmt.Errorf("missing token in header")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return refreshSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
