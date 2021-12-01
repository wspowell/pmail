package middleware

import (
	"github.com/wspowell/context"
	"github.com/wspowell/errors"
	"github.com/wspowell/spiderweb/httpheader"
	"github.com/wspowell/spiderweb/httpstatus"

	"github.com/wspowell/snailmail/resources/auth"
)

const (
	icMiddlewareUnknownAuthError = "middleware-auth-1"
	icMiddlewareInvalidToken     = "middleware-auth-2"
	icMiddlewareExpiredToken     = "middleware-auth-3"
	icMiddlewareTooEarly         = "middleware-auth-4"
)

var (
	ErrUserUnauthorized = errors.New("auth-1", "user unauthorized")
	ErrExpiredToken     = errors.New("auth-2", "expired token")
	ErrTooEarly         = errors.New("auth-3", "too early")
)

var (
	// nolint:gochecknoglobals // reason: spiderweb currently requires this to be a global
	JwtAuth auth.Jwt
)

type UserAuth struct {
	auth.SnailMailClaims
}

func (self *UserAuth) Authorization(ctx context.Context, peekHeader func(key string) []byte) (int, error) {
	authorizationHeaderValue := peekHeader(httpheader.Authorization)
	claims, err := JwtAuth.ValidateToken(string(authorizationHeaderValue))
	if err != nil {
		if errors.Is(err, auth.ErrTokenInvalid) {
			// FIXME: Need some way of setting Www-Authenticate header.
			return httpstatus.Unauthorized, errors.Convert(icMiddlewareInvalidToken, err, ErrUserUnauthorized)
		} else if errors.Is(err, auth.ErrTokenExpired) {
			return httpstatus.Unauthorized, errors.Convert(icMiddlewareExpiredToken, err, ErrExpiredToken)
		} else if errors.Is(err, auth.ErrTokenTooEarly) {
			return httpstatus.Unauthorized, errors.Convert(icMiddlewareTooEarly, err, ErrTooEarly)
		} else {
			return httpstatus.InternalServerError, errors.Convert(icMiddlewareUnknownAuthError, err, ErrUserUnauthorized)
		}
	}

	self.SnailMailClaims = *claims

	return httpstatus.OK, nil
}
