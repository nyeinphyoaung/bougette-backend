package helper

import (
	"bougette-backend/configs"
	"bougette-backend/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(user models.Users) (string, string, error) {
	userClaims := CustomClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)), // Refresh token valid for 1 hour
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	signedAccessToken, err := accessToken.SignedString([]byte(configs.Envs.JWT_SECRET))
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &CustomClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)), // Refresh token valid for 30 days
		},
	})
	signedRefreshToken, err := refreshToken.SignedString([]byte(configs.Envs.JWT_SECRET))
	if err != nil {
		return "", "", err

	}

	return signedAccessToken, signedRefreshToken, nil
}

func ParseAccessToken(tokenString string) (*CustomClaims, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(configs.Envs.JWT_SECRET), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, errors.New("token is malformed")
		} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, errors.New("access token is expired or not yet valid")
		} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, errors.New("invalid access token signature")
		}
		return nil, errors.New("failed to parse access token")
	}

	if !token.Valid {
		return nil, errors.New("invalid access token")
	}

	return claims, nil
}

func ParseRefreshToken(tokenString string) (*CustomClaims, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(configs.Envs.JWT_SECRET), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, errors.New("refresh token is malformed")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("refresh token is expired")
		} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, errors.New("invalid refresh token signature")
		}
		return nil, errors.New("failed to parse refresh token")
	}

	if !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	return claims, nil
}

func IsTokenExpired(claims *CustomClaims) bool {
	currentTime := jwt.NewNumericDate(time.Now())
	return claims.ExpiresAt.Time.Before(currentTime.Time)
}
