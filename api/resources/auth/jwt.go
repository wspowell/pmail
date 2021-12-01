package auth

import (
	"time"

	"github.com/wspowell/context"

	"github.com/golang-jwt/jwt/v4"
	"github.com/wspowell/errors"

	"github.com/wspowell/snailmail/resources/aws"
	"github.com/wspowell/snailmail/resources/models/user"
)

const (
	icTokenParseFailure = "auth-jwt-1"
	icTokenSignFailure  = "auth-jwt-2"
)

var (
	ErrTokenInvalid     = errors.New("jwt-1", "invalid token")
	ErrTokenExpired     = errors.New("jwt-2", "token expired")
	ErrTokenTooEarly    = errors.New("jwt-3", "token too early")
	ErrJwtSecretFailure = errors.New("jwt-4", "failed getting JWT signing key")
)

func GetSigningKey(ctx context.Context) ([]byte, error) {
	var jwtSignature jwtSigningKey
	if err := aws.GetSecret(ctx, &jwtSignature); err != nil {
		return nil, ErrJwtSecretFailure
	}

	return []byte(jwtSignature.Key), nil
}

type jwtSigningKey struct {
	Key string `json:"key"`
}

type SnailMailClaims struct {
	jwt.RegisteredClaims
	UserGuid          string `json:"userGuid"`
	Username          string `json:"username"`
	MailCarryCapacity uint32 `json:"mailCarryCapacity"`
	PineappleOnPizza  bool   `json:"pineappleOnPizza"`
}

type Group string

const (
	GroupUser  = Group("user")
	GroupAdmin = Group("admin")
)

func groupClaims(groups ...Group) jwt.ClaimStrings {
	claims := make(jwt.ClaimStrings, len(groups))
	for index := range groups {
		claims[index] = string(groups[index])
	}

	return claims
}

type Jwt struct {
	signingKey []byte
}

func NewJwt(signingKey []byte) Jwt {
	return Jwt{
		signingKey: signingKey,
	}
}

func (self Jwt) claims(authUser user.User, permissionGroups jwt.ClaimStrings) SnailMailClaims {
	return SnailMailClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "snailmail",
			Subject:   string(authUser.UserGuid),
			Audience:  permissionGroups,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
		UserGuid:          string(authUser.UserGuid),
		Username:          authUser.Username,
		MailCarryCapacity: authUser.MailCarryCapacity,
		PineappleOnPizza:  authUser.PineappleOnPizza,
	}
}

// ValidateToken
// Errors:
//   * ErrTokenInvalid
//   * ErrTokenExpired
//   * ErrTokenTooEarly
func (self Jwt) ValidateToken(tokenString string) (*SnailMailClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &SnailMailClaims{}, func(token *jwt.Token) (interface{}, error) {
		return self.signingKey, nil
	})
	if err != nil {
		return nil, errors.Propagate(icTokenParseFailure, ErrTokenInvalid)
	}

	if claims, ok := token.Claims.(*SnailMailClaims); ok && token.Valid {
		if time.Now().After(claims.RegisteredClaims.ExpiresAt.Time) {
			return nil, ErrTokenExpired
		}

		if time.Now().Before(claims.RegisteredClaims.NotBefore.Time) {
			return nil, ErrTokenTooEarly
		}

		return claims, nil
	}

	return nil, ErrTokenInvalid
}

func (self Jwt) UserToken(authUser user.User) (string, error) {
	jwtClaims := self.claims(authUser, groupClaims(GroupUser))

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	return self.signToken(token)
}

func (self Jwt) signToken(token *jwt.Token) (string, error) {
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(self.signingKey)
	if err != nil {
		return "", errors.Propagate(icTokenSignFailure, err)
	}

	return tokenString, nil
}
