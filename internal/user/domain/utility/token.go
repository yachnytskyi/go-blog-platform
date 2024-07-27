package utility

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/domain"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	signingMethod    = "RS256"
	userIDClaim      = "user_id"
	userRoleClaim    = "user_role"
	expirationClaim  = "exp"
	issuedAtClaim    = "iat"
	notBeforeClaim   = "nbf"
	unexpectedMethod = "unexpected method: %s"
)

// GenerateJWTToken generates a JWT token with the provided UserTokenPayload, using the given private key,
// and sets the token's expiration based on the specified token lifetime.
func GenerateJWTToken(logger model.Logger, location, privateKey string, tokenLifeTime time.Duration, userTokenPayload domain.UserTokenPayload) common.Result[string] {
	decodedPrivateKey := decodeBase64String(logger, location+".GenerateJWTToken", privateKey)
	if validator.IsError(decodedPrivateKey.Error) {
		return common.NewResultOnFailure[string](decodedPrivateKey.Error)
	}

	key := parsePrivateKey(logger, location+".GenerateJWTToken", decodedPrivateKey.Data)
	if validator.IsError(key.Error) {
		return common.NewResultOnFailure[string](key.Error)
	}

	now := time.Now().UTC()
	claims := generateClaims(tokenLifeTime, now, userTokenPayload)
	token := createSignedToken(logger, location+".GenerateJWTToken", key.Data, claims)
	if validator.IsError(token.Error) {
		return common.NewResultOnFailure[string](token.Error)
	}

	return common.NewResultOnSuccess[string](token.Data)
}

// ValidateJWTToken validates a JWT token using the provided public key and returns the claims
// extracted from the token if it's valid.
func ValidateJWTToken(logger model.Logger, location, token, publicKey string) common.Result[domain.UserTokenPayload] {
	decodedPublicKey := decodeBase64String(logger, location+".ValidateJWTToken", publicKey)
	if validator.IsError(decodedPublicKey.Error) {
		return common.NewResultOnFailure[domain.UserTokenPayload](decodedPublicKey.Error)
	}

	key := parsePublicKey(logger, location+".ValidateJWTToken", decodedPublicKey.Data)
	if validator.IsError(key.Error) {
		return common.NewResultOnFailure[domain.UserTokenPayload](key.Error)
	}

	parsedToken := parseToken(logger, location+".ValidateJWTToken", token, key.Data)
	if validator.IsError(parsedToken.Error) {
		return common.NewResultOnFailure[domain.UserTokenPayload](parsedToken.Error)
	}

	// Extract and validate the claims from the parsed token.
	claims, ok := parsedToken.Data.Claims.(jwt.MapClaims)
	if ok && parsedToken.Data.Valid {
		payload := domain.NewUserTokenPayload(
			fmt.Sprint(claims[userIDClaim]),
			fmt.Sprint(claims[userRoleClaim]),
		)

		return common.NewResultOnSuccess[domain.UserTokenPayload](payload)
	}

	invalidTokenError := domainError.NewInvalidTokenError(location+".ValidateJWTToken.Claims.ok", constants.InvalidTokenErrorMessage)
	logger.Error(invalidTokenError)
	return common.NewResultOnFailure[domain.UserTokenPayload](invalidTokenError)
}

// decodeBase64String decodes a base64-encoded string into a byte slice.
func decodeBase64String(logger model.Logger, location, base64String string) common.Result[[]byte] {
	decodedString, decodedStringError := base64.StdEncoding.DecodeString(base64String)
	if validator.IsError(decodedStringError) {
		internalError := domainError.NewInternalError(location+".decodeBase64String.StdEncoding.DecodeString", decodedStringError.Error())
		logger.Error(internalError)
		return common.NewResultOnFailure[[]byte](internalError)
	}

	return common.NewResultOnSuccess[[]byte](decodedString)
}

// generateClaims generates JWT claims with the specified token lifetime and UserTokenPayload.
func generateClaims(tokenLifeTime time.Duration, now time.Time, userTokenPayload domain.UserTokenPayload) jwt.MapClaims {
	return jwt.MapClaims{
		userIDClaim:     userTokenPayload.UserID,
		userRoleClaim:   userTokenPayload.Role,
		expirationClaim: now.Add(tokenLifeTime).Unix(),
		issuedAtClaim:   now.Unix(),
		notBeforeClaim:  now.Unix(),
	}
}

// createSignedToken creates a signed JWT token using the provided private key and claims.
func createSignedToken(logger model.Logger, location string, key *rsa.PrivateKey, claims jwt.MapClaims) common.Result[string] {
	token, tokenError := jwt.NewWithClaims(jwt.GetSigningMethod(signingMethod), claims).SignedString(key)
	if validator.IsError(tokenError) {
		internalError := domainError.NewInternalError(location+".createSignedToken.NewWithClaims.SignedString", tokenError.Error())
		logger.Error(internalError)
		return common.NewResultOnFailure[string](internalError)
	}

	return common.NewResultOnSuccess[string](token)
}

// parsePublicKey parses the RSA public key from the provided byte slice.
func parsePublicKey(logger model.Logger, location string, decodedPublicKey []byte) common.Result[*rsa.PublicKey] {
	key, keyError := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if validator.IsError(keyError) {
		internalError := domainError.NewInternalError(location+".parsePublicKey.ParseRSAPublicKeyFromPEM", keyError.Error())
		logger.Error(internalError)
		return common.NewResultOnFailure[*rsa.PublicKey](internalError)
	}

	return common.NewResultOnSuccess[*rsa.PublicKey](key)
}

// parsePrivateKey parses the RSA private key from the provided byte slice.
func parsePrivateKey(logger model.Logger, location string, decodedPrivateKey []byte) common.Result[*rsa.PrivateKey] {
	key, keyError := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if validator.IsError(keyError) {
		internalError := domainError.NewInternalError(location+".parsePrivateKey.ParseRSAPrivateKeyFromPEM", keyError.Error())
		logger.Error(internalError)
		return common.NewResultOnFailure[*rsa.PrivateKey](internalError)
	}

	return common.NewResultOnSuccess[*rsa.PrivateKey](key)
}

// parseToken parses and verifies the JWT token using the provided public key.
func parseToken(logger model.Logger, location, token string, key *rsa.PublicKey) common.Result[*jwt.Token] {
	parsedToken, parseError := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		_, ok := t.Method.(*jwt.SigningMethodRSA)
		if ok {
			return key, nil
		}

		internalError := domainError.NewInternalError(location+".parseToken.jwt.Parse.Ok", unexpectedMethod+" t.Header[alg]")
		logger.Error(internalError)
		return nil, internalError
	})

	if validator.IsError(parseError) {
		internalError := domainError.NewInternalError(location+"parseToken.jwt.Parse", parseError.Error())
		logger.Error(internalError)
		return common.NewResultOnFailure[*jwt.Token](internalError)
	}

	return common.NewResultOnSuccess[*jwt.Token](parsedToken)
}
