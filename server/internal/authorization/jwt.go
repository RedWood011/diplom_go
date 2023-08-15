// Package authorization пакет для авторизации
package authorization

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AtExpires    int64
	RtExpires    int64
}

func CreateToken(userID string, accessTokenLive time.Duration, refreshTokenLive time.Duration,
	accessTokenSecret string, refreshTokenSecret string) (*TokenDetails, error) {
	td := &TokenDetails{}

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userID
	atClaims["exp"] = time.Now().Add(accessTokenLive).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	accessToken, err := at.SignedString([]byte(accessTokenSecret))
	if err != nil {
		return nil, err
	}
	td.AccessToken = accessToken
	rtClaims := jwt.MapClaims{}
	rtClaims["user_id"] = userID
	rtClaims["exp"] = td.RtExpires

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(refreshTokenSecret))
	if err != nil {
		return nil, err
	}

	return td, nil
}

func verifyToken(ctx context.Context, accessSecret string) (*jwt.Token, error) {
	t := extractToken(ctx)
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signature is invalid : %v", token.Header["alg"])
		}
		return []byte(accessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(ctx context.Context, accessSecret string) (string, error) {
	token, err := verifyToken(ctx, accessSecret)
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", err
	}
	mapClaims := token.Claims.(jwt.MapClaims)
	userID := mapClaims["user_id"].(string)
	return userID, nil
}

func extractToken(ctx context.Context) string {
	token := metautils.ExtractIncoming(ctx).Get("authorization")
	array := strings.Split(token, " ")
	const typeAndToken = 2
	if len(array) == typeAndToken {
		return array[1]
	}
	return ""
}

func RefreshToken(refresh string, accessTokenLive time.Duration, refreshTokenLive time.Duration,
	accessSecret string, refreshSecret string) (*TokenDetails, error) {
	token, err := jwt.Parse(refresh, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signature is invalid: %v", token.Header["alg"])
		}
		return []byte("jdnfksdmfksd"), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userID := claims["user_id"].(string)
		ts, createErr := CreateToken(userID, accessTokenLive, refreshTokenLive,
			accessSecret, refreshSecret)
		if createErr != nil {
			return nil, err
		}
		return ts, nil
	}

	return nil, errors.New("refresh expired")
}
