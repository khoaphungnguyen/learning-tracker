package auth

import (
	"errors"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JwtWrapper struct {
	SecretKey         string // key used for signing the JWT token
	Issuer            string // Issuer of the JWT token
	ExpirationMinutes int64  // Number of minutes the JWT token will be valid for
	ExpirationHours   int64  // Expiration time of the JWT token in hours
}

// GenerateToken generates a jwt token
func (j *JwtWrapper) GenerateToken(id int) (signedToken string, err error) {
	claims := &jwt.StandardClaims{
		Audience:  strconv.Itoa(id),
		ExpiresAt: time.Now().Local().Add(time.Minute * time.Duration(j.ExpirationMinutes)).Unix(),
		Issuer:    j.Issuer,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err = token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return
	}
	return
}

// RefreshToken generates a refresh jwt token
func (j *JwtWrapper) RefreshToken(id int) (signedtoken string, err error) {
	claims := &jwt.StandardClaims{
		Audience:  strconv.Itoa(id),
		ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(j.ExpirationHours)).Unix(),
		Issuer:    j.Issuer,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedtoken, err = token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return
	}
	return
}

// ValidateToken validates the jwt token
func (j *JwtWrapper) ValidateToken(signedToken string) (claims *jwt.StandardClaims, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		err = errors.New("could not parse claims")
		return
	}
	if claims.ExpiresAt < time.Now().Unix() {
		err = errors.New("JWT is expired")
		return
	}
	return
}
