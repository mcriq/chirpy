package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	pw := []byte(password)
	dat, err := bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)
	if err != nil {
		log.Printf("error generating pw hash")
		return "", err
	}
	return string(dat), nil
}

func CheckPasswordHash(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	regClaims := jwt.RegisteredClaims{
		Issuer: "chirpy", 
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()), 
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)), 
		Subject: userID.String(),
	}
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, regClaims)
	tokenStr, err := newToken.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			return uuid.Nil, err
		}
		return userID, nil
	}
	
	return uuid.Nil, err
}

func GetBearerToken(headers http.Header) (string, error){
	if value, ok := headers["Authorization"]; ok {
		for _, val := range value {
			if str, ok := strings.CutPrefix(val, "Bearer "); ok {
				return str, nil
			}
		}
		return "", fmt.Errorf("no bearer token found")
	}
	return "", fmt.Errorf("no authorization found")
}

func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	
	return hex.EncodeToString(key), nil
}