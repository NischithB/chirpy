package utils

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJwt(key []byte, id int, isAccess bool) (string, error) {
	curTime := time.Now().UTC()
	var expiresIn time.Duration
	var issuer string

	if isAccess {
		issuer = "chirpy-access"
		expiresIn = time.Hour
	} else {
		issuer = "chirpy-refresh"
		expiresIn = time.Hour * time.Duration(24*60)
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(curTime),
			ExpiresAt: jwt.NewNumericDate(curTime.Add(expiresIn)),
			Subject:   fmt.Sprint(id),
		},
	)
	signedToken, err := token.SignedString(key)
	if err != nil {
		log.Printf("Failed to sign token: %s", err)
		return "", err
	}
	return signedToken, nil
}

func ValidateJwt(tokenString string, key []byte, isAccess bool) (int, error) {
	claim := jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(
		tokenString,
		&claim,
		func(t *jwt.Token) (interface{}, error) {
			return key, nil
		},
	)
	if err != nil {
		return -1, err
	}

	if isAccess && claim.Issuer != "chirpy-access" {
		return -1, jwt.ErrTokenInvalidIssuer
	}

	if !isAccess && claim.Issuer != "chirpy-refresh" {
		return -1, jwt.ErrTokenInvalidIssuer
	}

	id, err := strconv.Atoi(claim.Subject)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func ExtractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	split := strings.Split(authHeader, " ")
	if len(split) != 2 {
		return "", ErrTokenMissing
	}
	if split[1] == "" || split[1] == " " {
		return "", ErrTokenMissing
	}
	return split[1], nil
}
