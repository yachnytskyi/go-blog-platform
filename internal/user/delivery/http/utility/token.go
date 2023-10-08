package utility

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

func CreateToken(ttl time.Duration, payload interface{}, privateKey string) (string, error) {
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if validator.IsErrorNotNil(err) {
		return "", fmt.Errorf("could not decode key: %w", err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if validator.IsErrorNotNil(err) {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	now := time.Now().UTC()
	claims := make(jwt.MapClaims)
	claims["user_id"] = payload
	claims["exp"] = now.Add(ttl).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if validator.IsErrorNotNil(err) {
		return "", fmt.Errorf("create: sign token: %w", err)
	}
	return token, nil
}

func ValidateToken(token string, publicKey string) (interface{}, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if validator.IsErrorNotNil(err) {
		return nil, fmt.Errorf("could not decode: %w", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if validator.IsErrorNotNil(err) {
		return "", fmt.Errorf("validate: parse key: %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodRSA)
		if validator.IsBooleanNotTrue(ok) {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})
	if validator.IsErrorNotNil(err) {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if validator.IsBooleanNotTrue(ok) || validator.IsBooleanNotTrue(parsedToken.Valid) {
		return nil, fmt.Errorf("validate: invalid token")
	}
	return claims["user_id"], nil
}
