package security

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"good-template-go/conf"
	errors2 "good-template-go/pkg/errors"
	"good-template-go/pkg/utils"
	"time"
)

type ManagerClaims struct {
	ID        uuid.UUID `json:"id"`
	Role      string    `json:"role"`
	TokenType string    `json:"token_type"`
	jwt.RegisteredClaims
}

type ManagerToken struct {
	Token       string
	TokenType   string
	ExpiredTime time.Time
}

func GenerateManagerJWT(managerID uuid.UUID, role, tokenType string) (*ManagerToken, error) {
	expirationTokenTime := time.Now()
	secret := "secret"
	switch tokenType {
	case utils.AccessToken:
		secret = conf.GetConfig().JWTManagerAccessToken
		expirationTokenTime = time.Now().Add(time.Hour * 2)
	case utils.RefreshToken:
		secret = conf.GetConfig().JWTManagerRefreshToken
		expirationTokenTime = time.Now().Add(time.Hour * 10)
	default:
		appErr := errors2.FeAppError("Invalid token type", errors2.BadRequest)
		return nil, appErr
	}

	claims := ManagerClaims{
		ID:        managerID,
		Role:      role,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTokenTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	managerToken := ManagerToken{
		Token:       signedToken,
		TokenType:   tokenType,
		ExpiredTime: expirationTokenTime,
	}

	return &managerToken, nil
}

func ParseManagerJWT(tokenString string) (*ManagerClaims, error) {
	if tokenString == "" {
		return nil, errors.New("token is missing")
	}

	parser := jwt.Parser{}
	claims := &ManagerClaims{}
	_, _, err := parser.ParseUnverified(tokenString, claims)
	if err != nil {
		return nil, fmt.Errorf("failed to parse unverified token: %w", err)
	}

	if claims.Role != utils.Manager {
		return nil, errors.New("invalid role: only manager tokens are allowed")
	}

	var secret string
	if claims.TokenType == utils.AccessToken {
		secret = conf.GetConfig().JWTManagerAccessToken
	} else {
		return nil, errors.New("invalid token type: only access tokens are allowed")
	}

	token, err := jwt.ParseWithClaims(tokenString, &ManagerClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("manager token expired")
		}
		return nil, err
	}

	claims, ok := token.Claims.(*ManagerClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid manager token")
	}

	return claims, nil
}
