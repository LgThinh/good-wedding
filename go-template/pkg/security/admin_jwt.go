package security

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"good-template-go/conf"
	errors2 "good-template-go/pkg/errors"
	"good-template-go/pkg/utils"
	"time"
)

type AdminClaims struct {
	ID        uuid.UUID `json:"id"`
	Role      string    `json:"role"`
	TokenType string    `json:"token_type"`
	jwt.RegisteredClaims
}

type AdminToken struct {
	Token       string    `json:"token"`
	TokenType   string    `json:"token_type"`
	ExpiredTime time.Time `json:"expired_time"`
}

func GenerateAdminJWT(adminID uuid.UUID, role, tokenType string) (*AdminToken, error) {
	expirationTokenTime := time.Now()
	secret := "secret"
	switch tokenType {
	case utils.AccessToken:
		secret = conf.GetConfig().JWTAdminAccessToken
		expirationTokenTime = time.Now().Add(time.Hour * 2)
	case utils.RefreshToken:
		secret = conf.GetConfig().JWTAdminRefreshToken
		expirationTokenTime = time.Now().Add(time.Hour * 10)
	default:
		appErr := errors2.FeAppError("Invalid token type", errors2.BadRequest)
		return nil, appErr
	}

	claims := AdminClaims{
		ID:        adminID,
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

	adminToken := AdminToken{
		Token:       signedToken,
		TokenType:   tokenType,
		ExpiredTime: expirationTokenTime,
	}

	return &adminToken, nil
}

func ParseAdminJWT(tokenString string) (*AdminClaims, error) {
	if tokenString == "" {
		return nil, errors.New("token is missing")
	}

	parser := jwt.Parser{}
	claims := &AdminClaims{}
	_, _, err := parser.ParseUnverified(tokenString, claims)
	if err != nil {
		return nil, fmt.Errorf("failed to parse unverified token: %w", err)
	}

	if claims.Role != utils.Admin {
		return nil, errors.New("invalid role: only admin tokens are allowed")
	}

	var secret string
	if claims.TokenType == utils.AccessToken {
		secret = conf.GetConfig().JWTAdminAccessToken
	} else {
		return nil, errors.New("invalid token type: only access tokens are allowed")
	}

	token, err := jwt.ParseWithClaims(tokenString, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("admin token expired")
		}
		return nil, err
	}

	claims, ok := token.Claims.(*AdminClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid admin token")
	}

	return claims, nil
}
