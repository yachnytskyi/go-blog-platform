package utility

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domainModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/domain"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logger "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logger"
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
func GenerateJWTToken(location, privateKey string, tokenLifeTime time.Duration, userTokenPayload domainModel.UserTokenPayload) commonModel.Result[string] {
	decodedPrivateKey := decodeBase64String(location+".GenerateJWTToken", privateKey)
	if validator.IsError(decodedPrivateKey.Error) {
		return commonModel.NewResultOnFailure[string](decodedPrivateKey.Error)
	}

	key := parsePrivateKey(location+".GenerateJWTToken", decodedPrivateKey.Data)
	if validator.IsError(key.Error) {
		return commonModel.NewResultOnFailure[string](key.Error)
	}

	now := time.Now().UTC()
	claims := generateClaims(tokenLifeTime, now, userTokenPayload)
	token := createSignedToken(location+".GenerateJWTToken", key.Data, claims)
	if validator.IsError(token.Error) {
		return commonModel.NewResultOnFailure[string](token.Error)
	}

	return commonModel.NewResultOnSuccess[string](token.Data)
}

// ValidateJWTToken validates a JWT token using the provided public key and returns the claims
// extracted from the token if it's valid.
func ValidateJWTToken(location, token, publicKey string) commonModel.Result[domainModel.UserTokenPayload] {
	decodedPublicKey := decodeBase64String(location+".ValidateJWTToken", publicKey)
	if validator.IsError(decodedPublicKey.Error) {
		return commonModel.NewResultOnFailure[domainModel.UserTokenPayload](decodedPublicKey.Error)
	}

	key := parsePublicKey(location+".ValidateJWTToken", decodedPublicKey.Data)
	if validator.IsError(key.Error) {
		return commonModel.NewResultOnFailure[domainModel.UserTokenPayload](key.Error)
	}

	parsedToken := parseToken(location+".ValidateJWTToken", token, key.Data)
	if validator.IsError(parsedToken.Error) {
		return commonModel.NewResultOnFailure[domainModel.UserTokenPayload](parsedToken.Error)
	}

	// Extract and validate the claims from the parsed token.
	claims, ok := parsedToken.Data.Claims.(jwt.MapClaims)
	if ok && parsedToken.Data.Valid {
		payload := domainModel.NewUserTokenPayload(
			fmt.Sprint(claims[userIDClaim]),
			fmt.Sprint(claims[userRoleClaim]),
		)

		return commonModel.NewResultOnSuccess[domainModel.UserTokenPayload](payload)
	}

	invalidTokenError := domainError.NewInvalidTokenError(location+".ValidateJWTToken.Claims.ok", constants.InvalidTokenErrorMessage)
	logger.Logger(invalidTokenError)
	return commonModel.NewResultOnFailure[domainModel.UserTokenPayload](invalidTokenError)
}

// decodeBase64String decodes a base64-encoded string into a byte slice.
func decodeBase64String(location, base64String string) commonModel.Result[[]byte] {
	decodedPrivateKey, decodeStringError := base64.StdEncoding.DecodeString(base64String)
	if validator.IsError(decodeStringError) {
		internalError := domainError.NewInternalError(location+".decodeBase64String.StdEncoding.DecodeString", decodeStringError.Error())
		logger.Logger(internalError)
		return commonModel.NewResultOnFailure[[]byte](internalError)
	}

	return commonModel.NewResultOnSuccess[[]byte](decodedPrivateKey)
}

// generateClaims generates JWT claims with the specified token lifetime and UserTokenPayload.
func generateClaims(tokenLifeTime time.Duration, now time.Time, userTokenPayload domainModel.UserTokenPayload) jwt.MapClaims {
	return jwt.MapClaims{
		userIDClaim:     userTokenPayload.UserID,
		userRoleClaim:   userTokenPayload.Role,
		expirationClaim: now.Add(tokenLifeTime).Unix(),
		issuedAtClaim:   now.Unix(),
		notBeforeClaim:  now.Unix(),
	}
}

// createSignedToken creates a signed JWT token using the provided private key and claims.
func createSignedToken(location string, key *rsa.PrivateKey, claims jwt.MapClaims) commonModel.Result[string] {
	token, createSignedTokenError := jwt.NewWithClaims(jwt.GetSigningMethod(signingMethod), claims).SignedString(key)
	if validator.IsError(createSignedTokenError) {
		internalError := domainError.NewInternalError(location+".createSignedToken.NewWithClaims", createSignedTokenError.Error())
		logger.Logger(internalError)
		return commonModel.NewResultOnFailure[string](internalError)
	}

	return commonModel.NewResultOnSuccess[string](token)
}

// parsePublicKey parses the RSA public key from the provided byte slice.
func parsePublicKey(location string, decodedPublicKey []byte) commonModel.Result[*rsa.PublicKey] {
	key, parsePublicKeyError := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if validator.IsError(parsePublicKeyError) {
		internalError := domainError.NewInternalError(location+".parsePublicKey.ParseRSAPublicKeyFromPEM", parsePublicKeyError.Error())
		logger.Logger(internalError)
		return commonModel.NewResultOnFailure[*rsa.PublicKey](internalError)
	}

	return commonModel.NewResultOnSuccess[*rsa.PublicKey](key)
}

// parsePrivateKey parses the RSA private key from the provided byte slice.
func parsePrivateKey(location string, decodedPrivateKey []byte) commonModel.Result[*rsa.PrivateKey] {
	key, parsePrivateKeyError := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if validator.IsError(parsePrivateKeyError) {
		internalError := domainError.NewInternalError(location+".parsePrivateKey.ParseRSAPrivateKeyFromPEM", parsePrivateKeyError.Error())
		logger.Logger(internalError)
		return commonModel.NewResultOnFailure[*rsa.PrivateKey](internalError)
	}

	return commonModel.NewResultOnSuccess[*rsa.PrivateKey](key)
}

// parseToken parses and verifies the JWT token using the provided public key.
func parseToken(location, token string, key *rsa.PublicKey) commonModel.Result[*jwt.Token] {
	parsedToken, parseError := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		_, ok := t.Method.(*jwt.SigningMethodRSA)
		if ok {
			return key, nil
		}

		internalError := domainError.NewInternalError(location+".parseToken.jwt.Parse.Ok", unexpectedMethod+" t.Header[alg]")
		logger.Logger(internalError)
		return nil, internalError
	})

	if validator.IsError(parseError) {
		internalError := domainError.NewInternalError(location+"parseToken.jwt.Parse", parseError.Error())
		logger.Logger(internalError)
		return commonModel.NewResultOnFailure[*jwt.Token](internalError)
	}

	return commonModel.NewResultOnSuccess[*jwt.Token](parsedToken)
}
