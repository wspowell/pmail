package middleware

import (
	"github.com/wspowell/snailmail/resources/auth"

	"github.com/wspowell/context"
	"github.com/wspowell/errors"
	"github.com/wspowell/spiderweb/httpheader"
	"github.com/wspowell/spiderweb/httpstatus"
)

const (
	icMiddlewareInvalidToken     = "middleware-auth-1"
	icMiddlewareExpiredToken     = "middleware-auth-2"
	icMiddlewareUnknownAuthError = "middleware-auth-3"
)

var (
	ErrUserUnauthorized = errors.New("auth-1", "user unauthorized")
	ErrExpiredToken     = errors.New("auth-2", "expired token")
	ErrTooEarly         = errors.New("auth-3", "too early")
)

var (
	JwtAuth auth.Jwt
)

type UserAuth struct {
	auth.SnailMailClaims
}

func (self *UserAuth) Authorization(ctx context.Context, PeekHeader func(key string) []byte) (int, error) {
	authorizationHeaderValue := PeekHeader(httpheader.Authorization)
	claims, err := JwtAuth.ValidateToken(string(authorizationHeaderValue))
	if err != nil {
		if errors.Is(err, auth.ErrTokenInvalid) {
			// FIXME: Need some way of setting Www-Authenticate header.
			return httpstatus.Unauthorized, errors.Convert(icMiddlewareInvalidToken, err, ErrUserUnauthorized)
		} else if errors.Is(err, auth.ErrTokenExpired) {
			return httpstatus.Unauthorized, errors.Convert(icMiddlewareExpiredToken, err, ErrExpiredToken)
		} else if errors.Is(err, auth.ErrTokenTooEarly) {
			return httpstatus.Unauthorized, errors.Convert(icMiddlewareExpiredToken, err, ErrTooEarly)
		} else {
			return httpstatus.InternalServerError, errors.Convert(icMiddlewareUnknownAuthError, err, ErrUserUnauthorized)
		}
	}

	self.SnailMailClaims = *claims

	return httpstatus.OK, err
}
