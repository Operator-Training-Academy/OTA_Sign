package app

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidToken = errors.New("invalid launch token")
	ErrExpiredToken = errors.New("expired launch token")
)

func SignLaunchClaims(claims LaunchClaims, secret string) (string, error) {
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	return encodedPayload + "." + signLaunchPayload(encodedPayload, secret), nil
}

func VerifyLaunchToken(token string, secret string, now time.Time) (LaunchClaims, error) {
	var claims LaunchClaims
	parts := strings.Split(token, ".")
	if len(parts) != 2 || !hmac.Equal([]byte(signLaunchPayload(parts[0], secret)), []byte(parts[1])) {
		return claims, ErrInvalidToken
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil || json.Unmarshal(payload, &claims) != nil {
		return claims, ErrInvalidToken
	}
	if claims.ExpiresAt.Before(now) {
		return claims, ErrExpiredToken
	}
	return claims, nil
}

func signLaunchPayload(encodedPayload string, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(encodedPayload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
