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
	signingMethod    = "RS256"
	user_id          = "user_id"
	exp              = "exp"
	iat              = "iat"
	nbf              = "nbf"
	location         = "internal.user.delivery.http.utility."
	unexpectedMethod = "unexpected method: %s"
	invalidToken     = "validate: invalid token"
)

func CreateToken(ttl time.Duration, payload interface{}, privateKey string) (string, error) {
	decodedPrivateKey, decodeStringError := base64.StdEncoding.DecodeString(privateKey)
	if validator.IsErrorNotNil(decodeStringError) {
		decodeStringInternalError := domainError.NewInternalError(location+"CreateToken.StdEncoding.DecodeString", decodeStringError.Error())
		logging.Logger(decodeStringInternalError)
		return "", decodeStringInternalError
	}
	key, parsePrivateKeyError := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if validator.IsErrorNotNil(parsePrivateKeyError) {
		parsePrivateKeyInternalError := domainError.NewInternalError(location+"CreateToken.ParseRSAPrivateKeyFromPEM", parsePrivateKeyError.Error())
		logging.Logger(parsePrivateKeyInternalError)
		return "", parsePrivateKeyInternalError
	}
	now := time.Now().UTC()
	claims := jwt.MapClaims{
		user_id: payload,
		exp:     now.Add(ttl).Unix(),
		iat:     now.Unix(),
		nbf:     now.Unix(),
	}
	token, newWithClaimsError := jwt.NewWithClaims(jwt.GetSigningMethod(signingMethod), claims).SignedString(key)
	if validator.IsErrorNotNil(newWithClaimsError) {
		newWithClaimsInternalError := domainError.NewInternalError(location+"CreateToken.NewWithClaims", newWithClaimsError.Error())
		logging.Logger(newWithClaimsInternalError)
		return "", newWithClaimsInternalError
	}
	return token, nil
}

func ValidateToken(token string, publicKey string) (interface{}, error) {
	decodedPublicKey, decodeStringError := base64.StdEncoding.DecodeString(publicKey)
	if validator.IsErrorNotNil(decodeStringError) {
		decodeStringInternalError := domainError.NewInternalError(location+"ValidateToken.DecodeString", decodeStringError.Error())
		logging.Logger(decodeStringInternalError)
		return nil, decodeStringInternalError
	}
	key, parseError := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if validator.IsErrorNotNil(parseError) {
		parseInternalError := domainError.NewInternalError(location+"ValidateToken.DecodeString", parseError.Error())
		logging.Logger(parseInternalError)
		return nil, parseInternalError
	}
	parsedToken, parseError := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodRSA)
		if validator.IsBooleanNotTrue(ok) {
			internalError := domainError.NewInternalError(location+"ValidateToken.jwt.Parse.NotOk", unexpectedMethod+" t.Header[alg]")
			logging.Logger(internalError)
			return nil, internalError
		}
		return key, nil
	})
	if validator.IsErrorNotNil(parseError) {
		parseInternalError := domainError.NewInternalError(location+"ValidateToken.jwt.Parse", parseError.Error())
		logging.Logger(parseInternalError)
		return nil, parseInternalError
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if validator.IsBooleanNotTrue(ok) || validator.IsBooleanNotTrue(parsedToken.Valid) {
		errorMessage := domainError.NewInternalError(location+"parsedToken.Claims.NotOk", invalidToken)
		logging.Logger(errorMessage)
		return nil, errorMessage
	}
	return claims[user_id], nil
}
