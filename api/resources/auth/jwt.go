package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/wspowell/errors"
	"github.com/wspowell/snailmail/resources/models/user"
)

const (
	icTokenSignFailure  = "auth-jwt-1"
	icTokenParseFailure = "auth-jwt-2"
	icInvalidToken      = "auth-jwt-3"
	icTokenExpired      = "auth-jwt-4"
	icTokenEarly        = "auth-jwt-5"
)

var (
	ErrTokenInvalid  = errors.New(icInvalidToken, "invalid token")
	ErrTokenExpired  = errors.New(icTokenExpired, "token expired")
	ErrTokenTooEarly = errors.New(icTokenEarly, "token too early")
)

var (
	// TODO: Retrieve secret key from somewhere safe.
	signingKey = []byte("tempkey")
)

func GetSigningKey() ([]byte, error) {
	return signingKey, nil
}

type SnailMailClaims struct {
	jwt.RegisteredClaims
	UserGuid          string `json:"user_guid"`
	Username          string `json:"username"`
	MailCarryCapacity uint32 `json:"mail_carry_capacity"`
	PineappleOnPizza  bool   `json:"pineapple_on_pizza"`
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
