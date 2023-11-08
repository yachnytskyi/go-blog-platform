package utility

import (
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	signingMethod       = "RS256"
	userIDClaim         = "user_id"
	expirationClaim     = "exp"
	issuedAtClaim       = "iat"
	notBeforeClaim      = "nbf"
	location            = "internal.user.domain.utility."
	unexpectedMethod    = "unexpected method: %s"
	invalidTokenMessage = "validate: invalid token"
)

func GenerateJWTToken(tokenLifeTime time.Duration, payload any, privateKey string) (string, error) {
	decodedPrivateKey, decodeStringError := base64.StdEncoding.DecodeString(privateKey)
	if validator.IsErrorNotNil(decodeStringError) {
		internalError := domainError.NewInternalError(location+"CreateToken.StdEncoding.DecodeString", decodeStringError.Error())
		logging.Logger(internalError)
		return "", internalError
	}
	key, parsePrivateKeyError := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if validator.IsErrorNotNil(parsePrivateKeyError) {
		internalError := domainError.NewInternalError(location+"CreateToken.ParseRSAPrivateKeyFromPEM", parsePrivateKeyError.Error())
		logging.Logger(internalError)
		return "", internalError
	}
	now := time.Now().UTC()
	claims := jwt.MapClaims{
		userIDClaim:     payload,
		expirationClaim: now.Add(tokenLifeTime).Unix(),
		issuedAtClaim:   now.Unix(),
		notBeforeClaim:  now.Unix(),
	}
	token, newWithClaimsError := jwt.NewWithClaims(jwt.GetSigningMethod(signingMethod), claims).SignedString(key)
	if validator.IsErrorNotNil(newWithClaimsError) {
		internalError := domainError.NewInternalError(location+"CreateToken.NewWithClaims", newWithClaimsError.Error())
		logging.Logger(internalError)
		return "", internalError
	}
	return token, nil
}

func ValidateJWTToken(token string, publicKey string) (any, error) {
	decodedPublicKey, decodeStringError := base64.StdEncoding.DecodeString(publicKey)
	if validator.IsErrorNotNil(decodeStringError) {
		internalError := domainError.NewInternalError(location+"ValidateToken.DecodeString", decodeStringError.Error())
		logging.Logger(internalError)
		return nil, internalError
	}
	key, parseError := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if validator.IsErrorNotNil(parseError) {
		internalError := domainError.NewInternalError(location+"ValidateToken.DecodeString", parseError.Error())
		logging.Logger(internalError)
		return nil, internalError
	}
	parsedToken, parseError := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		_, ok := t.Method.(*jwt.SigningMethodRSA)
		if validator.IsBooleanNotTrue(ok) {
			internalError := domainError.NewInternalError(location+"ValidateToken.jwt.Parse.NotOk", unexpectedMethod+" t.Header[alg]")
			logging.Logger(internalError)
			return nil, internalError
		}
		return key, nil
	})
	if validator.IsErrorNotNil(parseError) {
		internalError := domainError.NewInternalError(location+"ValidateToken.jwt.Parse", parseError.Error())
		logging.Logger(internalError)
		return nil, internalError
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if validator.IsBooleanNotTrue(ok) || validator.IsBooleanNotTrue(parsedToken.Valid) {
		errorMessage := domainError.NewInternalError(location+"parsedToken.Claims.NotOk", invalidTokenMessage)
		logging.Logger(errorMessage)
		return nil, errorMessage
	}
	return claims[userIDClaim], nil
}
