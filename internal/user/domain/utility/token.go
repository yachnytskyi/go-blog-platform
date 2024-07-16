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
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
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
//
// Parameters:
// - location: A string representing the location for logging purposes.
// - privateKey: A base64-encoded string representing the RSA private key used for signing the JWT token.
// - tokenLifeTime: A duration representing the lifetime of the token.
// - userTokenPayload: A UserTokenPayload struct containing the user's ID and role.
//
// Returns:
// - commonModel.Result[string]: The result containing either the generated JWT token string or an error.
func GenerateJWTToken(location, privateKey string, tokenLifeTime time.Duration, userTokenPayload domainModel.UserTokenPayload) commonModel.Result[string] {
	// Decode the private key from base64-encoded string.
	decodedPrivateKey := decodeBase64String(location+".GenerateJWTToken", privateKey)
	if validator.IsError(decodedPrivateKey.Error) {
		return commonModel.NewResultOnFailure[string](decodedPrivateKey.Error)
	}

	// Parse the private key for signing.
	key := parsePrivateKey(location+".GenerateJWTToken", decodedPrivateKey.Data)
	if validator.IsError(key.Error) {
		return commonModel.NewResultOnFailure[string](key.Error)
	}

	// Generate claims for the JWT token.
	// Create the signed token using the private key and claims.
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
//
// Parameters:
// - location: A string representing the location for logging purposes.
// - token: A string representing the JWT token to be validated.
// - publicKey: A base64-encoded string representing the RSA public key used for verifying the JWT token.
//
// Returns:
// - commonModel.Result[domainModel.UserTokenPayload]: The result containing either the extracted UserTokenPayload if the token is valid, or an error.
func ValidateJWTToken(location, token, publicKey string) commonModel.Result[domainModel.UserTokenPayload] {
	// Decode the public key from a base64-encoded string.
	decodedPublicKey := decodeBase64String(location+".ValidateJWTToken", publicKey)
	if validator.IsError(decodedPublicKey.Error) {
		return commonModel.NewResultOnFailure[domainModel.UserTokenPayload](decodedPublicKey.Error)
	}

	// Parse the public key for verification.
	key := parsePublicKey(location+".ValidateJWTToken", decodedPublicKey.Data)
	if validator.IsError(key.Error) {
		return commonModel.NewResultOnFailure[domainModel.UserTokenPayload](key.Error)
	}

	// Parse and verify the token using the public key.
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
	logging.Logger(invalidTokenError)
	return commonModel.NewResultOnFailure[domainModel.UserTokenPayload](invalidTokenError)
}

// decodeBase64String decodes a base64-encoded string into a byte slice.
// It performs the following steps:
// 1. Attempts to decode the provided base64 string.
// 2. Checks for errors during decoding and logs them if any.
// 3. Returns the decoded byte slice wrapped in a commonModel.Result.
//
// Parameters:
// - location (string): A string representing the location for logging purposes.
// - base64String (string): A base64-encoded string to be decoded.
//
// Returns:
// - commonModel.Result[[]byte]: The result containing either the decoded byte slice or an error.
func decodeBase64String(location, base64String string) commonModel.Result[[]byte] {
	decodedPrivateKey, decodeStringError := base64.StdEncoding.DecodeString(base64String)
	if validator.IsError(decodeStringError) {
		internalError := domainError.NewInternalError(location+".decodeBase64String.StdEncoding.DecodeString", decodeStringError.Error())
		logging.Logger(internalError)
		return commonModel.NewResultOnFailure[[]byte](internalError)
	}

	return commonModel.NewResultOnSuccess[[]byte](decodedPrivateKey)
}

// generateClaims generates JWT claims with the specified token lifetime and UserTokenPayload.
// It sets standard claims like userID, userRole, expiration, issuedAt, and notBefore.
//
// Parameters:
// - tokenLifeTime (time.Duration): A duration representing the lifetime of the token.
// - now (time.Time): The current time used for setting the issuedAt and notBefore claims.
// - userTokenPayload (domainModel.UserTokenPayload): A struct containing the user's ID and role.
//
// Returns:
// - jwt.MapClaims: A MapClaims object containing the generated claims.
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
// It performs the following steps:
// 1. Creates a new JWT token with the provided claims.
// 2. Attempts to sign the token using the provided private key.
// 3. Checks for errors during signing and logs them if any.
// 4. Returns the signed token string wrapped in a commonModel.Result.
//
// Parameters:
// - location (string): A string representing the location for logging purposes.
// - key (*rsa.PrivateKey): A pointer to the RSA private key used for signing the JWT token.
// - claims (jwt.MapClaims): A MapClaims object containing the claims to be included in the JWT token.
//
// Returns:
// - commonModel.Result[string]: The result containing either the signed JWT token string or an error.
func createSignedToken(location string, key *rsa.PrivateKey, claims jwt.MapClaims) commonModel.Result[string] {
	token, createSignedTokenError := jwt.NewWithClaims(jwt.GetSigningMethod(signingMethod), claims).SignedString(key)
	if validator.IsError(createSignedTokenError) {
		internalError := domainError.NewInternalError(location+".createSignedToken.NewWithClaims", createSignedTokenError.Error())
		logging.Logger(internalError)
		return commonModel.NewResultOnFailure[string](internalError)
	}

	return commonModel.NewResultOnSuccess[string](token)
}

// parsePublicKey parses the RSA public key from the provided byte slice.
// It performs the following steps:
// 1. Attempts to parse the provided byte slice into an RSA public key.
// 2. Checks for errors during parsing and logs them if any.
// 3. Returns the parsed RSA public key wrapped in a commonModel.Result.
//
// Parameters:
// - location (string): A string representing the location for logging purposes.
// - decodedPublicKey ([]byte): A byte slice containing the decoded RSA public key.
//
// Returns:
// - commonModel.Result[*rsa.PublicKey]: The result containing either the parsed RSA public key or an error.
func parsePublicKey(location string, decodedPublicKey []byte) commonModel.Result[*rsa.PublicKey] {
	key, parsePublicKeyError := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if validator.IsError(parsePublicKeyError) {
		internalError := domainError.NewInternalError(location+".parsePublicKey.ParseRSAPublicKeyFromPEM", parsePublicKeyError.Error())
		logging.Logger(internalError)
		return commonModel.NewResultOnFailure[*rsa.PublicKey](internalError)
	}

	return commonModel.NewResultOnSuccess[*rsa.PublicKey](key)
}

// parsePrivateKey parses the RSA private key from the provided byte slice.
// It performs the following steps:
// 1. Attempts to parse the provided byte slice into an RSA private key.
// 2. Checks for errors during parsing and logs them if any.
// 3. Returns the parsed RSA private key wrapped in a commonModel.Result.
//
// Parameters:
// - location (string): A string representing the location for logging purposes.
// - decodedPrivateKey ([]byte): A byte slice containing the decoded RSA private key.
//
// Returns:
// - commonModel.Result[*rsa.PrivateKey]: The result containing either the parsed RSA private key or an error.
func parsePrivateKey(location string, decodedPrivateKey []byte) commonModel.Result[*rsa.PrivateKey] {
	key, parsePrivateKeyError := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if validator.IsError(parsePrivateKeyError) {
		internalError := domainError.NewInternalError(location+".parsePrivateKey.ParseRSAPrivateKeyFromPEM", parsePrivateKeyError.Error())
		logging.Logger(internalError)
		return commonModel.NewResultOnFailure[*rsa.PrivateKey](internalError)
	}

	return commonModel.NewResultOnSuccess[*rsa.PrivateKey](key)
}

// parseToken parses and verifies the JWT token using the provided public key.
// It performs the following steps:
// 1. Attempts to parse the provided JWT token.
// 2. Verifies the token using the provided public key.
// 3. Checks for errors during parsing and verification and logs them if any.
// 4. Returns the parsed token wrapped in a commonModel.Result.
//
// Parameters:
// - location (string): A string representing the location for logging purposes.
// - token (string): A string representing the JWT token to be parsed and verified.
// - key (*rsa.PublicKey): A pointer to the RSA public key used for verifying the JWT token.
//
// Returns:
// - commonModel.Result[*jwt.Token]: The result containing either the parsed JWT token or an error.
func parseToken(location, token string, key *rsa.PublicKey) commonModel.Result[*jwt.Token] {
	parsedToken, parseError := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		_, ok := t.Method.(*jwt.SigningMethodRSA)
		if ok {
			return key, nil
		}

		internalError := domainError.NewInternalError(location+".parseToken.jwt.Parse.Ok", unexpectedMethod+" t.Header[alg]")
		logging.Logger(internalError)
		return nil, internalError
	})

	if validator.IsError(parseError) {
		internalError := domainError.NewInternalError(location+"parseToken.jwt.Parse", parseError.Error())
		logging.Logger(internalError)
		return commonModel.NewResultOnFailure[*jwt.Token](internalError)
	}

	return commonModel.NewResultOnSuccess[*jwt.Token](parsedToken)
}
