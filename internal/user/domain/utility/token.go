package utility

import (
	"crypto/rsa"
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

// GenerateJWTToken generates a JWT token with the provided payload, using the given private key,
// and sets the token's expiration based on the specified token lifetime.
func GenerateJWTToken(tokenLifeTime time.Duration, payload any, privateKey string) (string, error) {
	// Decode the private key from base64-encoded string.
	decodedPrivateKey, decodeStringError := decodeString(privateKey)
	if validator.IsErrorNotNil(decodeStringError) {
		return "", decodeStringError
	}

	// Parse the private key for signing.
	key, parsePrivateKeyError := parsePrivateKey(decodedPrivateKey)
	if validator.IsErrorNotNil(parsePrivateKeyError) {
		return "", parsePrivateKeyError
	}

	// Generate claims for the JWT token.
	// Create the signed token using the private key and claims.
	now := time.Now().UTC()
	claims := generateClaims(tokenLifeTime, now, payload)
	token, newWithClaimsError := createSignedToken(key, claims)
	if validator.IsErrorNotNil(newWithClaimsError) {
		return "", newWithClaimsError
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

// decodeString decodes a base64-encoded string into a byte slice.
func decodeString(decodedString string) ([]byte, error) {
	decodedPrivateKey, decodeStringError := base64.StdEncoding.DecodeString(decodedString)
	if validator.IsErrorNotNil(decodeStringError) {
		internalError := domainError.NewInternalError(location+"decodeString.StdEncoding.DecodeString", decodeStringError.Error())
		logging.Logger(internalError)
		return []byte{}, internalError
	}
	return decodedPrivateKey, nil
}

// parsePrivateKey parses the RSA private key from the provided byte slice.
func parsePrivateKey(decodedPrivateKey []byte) (*rsa.PrivateKey, error) {
	key, parsePrivateKeyError := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if validator.IsErrorNotNil(parsePrivateKeyError) {
		internalError := domainError.NewInternalError(location+"parsePrivateKey.ParseRSAPrivateKeyFromPEM", parsePrivateKeyError.Error())
		logging.Logger(internalError)
		return nil, internalError
	}
	return key, nil
}

// generateClaims generates JWT claims with the specified token lifetime and payload.
func generateClaims(tokenLifeTime time.Duration, now time.Time, payload any) jwt.MapClaims {
	return jwt.MapClaims{
		userIDClaim:     payload,
		expirationClaim: now.Add(tokenLifeTime).Unix(),
		issuedAtClaim:   now.Unix(),
		notBeforeClaim:  now.Unix(),
	}
}

// createSignedToken creates a signed JWT token using the provided private key and claims.
func createSignedToken(key *rsa.PrivateKey, claims jwt.MapClaims) (string, error) {
	token, newWithClaimsError := jwt.NewWithClaims(jwt.GetSigningMethod(signingMethod), claims).SignedString(key)
	if newWithClaimsError != nil {
		internalError := domainError.NewInternalError(location+"createSignedToken.NewWithClaims", newWithClaimsError.Error())
		logging.Logger(internalError)
		return "", internalError
	}
	return token, nil
}
