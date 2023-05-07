package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

// JWTService is a service for JWT operations
type JWTService struct {
	signingKey      string
	tokenExpiration int // The expiration time in minutes
}

func NewJWTService(signingKey string, tokenExpiration int) *JWTService {
	return &JWTService{
		signingKey:      signingKey,
		tokenExpiration: tokenExpiration,
	}
}

func (jwtSvc *JWTService) ValidateToken(authToken string) (*jwt.Token, error) {
	return jwt.Parse(authToken, func(t_ *jwt.Token) (interface{}, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("bad signed method received")
		}
		return []byte(jwtSvc.signingKey), nil
	})
}

func (jwtSvc *JWTService) GenerateToken(identity Identity) (string, error) {
	expirationTime := time.Now().Add(time.Duration(jwtSvc.tokenExpiration) * time.Second)
	claims := &Claims{
		ID: identity.GetID(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSvc.signingKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
