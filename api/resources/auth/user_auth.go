package auth

import (
	"github.com/wspowell/errors"
	"github.com/wspowell/spiderweb/httpheader"
	"github.com/wspowell/spiderweb/httpstatus"
	"github.com/wspowell/spiderweb/httptrip"
)

var (
	ErrUserUnauthorized = errors.New("user unauthorized")
	ErrExpiredToken     = errors.New("expired token")
	ErrTooEarly         = errors.New("too early")
)

type User struct {
	JwtAuth        Jwt
	AuthorizedUser SnailMailClaims
}

func (self *User) Authorize(reqRes httptrip.RoundTripper) error {
	authorizationHeaderValue := reqRes.PeekRequestHeader(httpheader.Authorization)
	claims, err := self.JwtAuth.ValidateToken(string(authorizationHeaderValue))
	if err != nil {
		if errors.Is(err, ErrTokenInvalid) {
			reqRes.SetResponseHeader(httpheader.WwwAuthenticate, "JWT Token")
			reqRes.SetStatusCode(httpstatus.Unauthorized)
			return errors.Wrap(err, ErrUserUnauthorized)
		} else if errors.Is(err, ErrTokenExpired) {
			reqRes.SetStatusCode(httpstatus.Unauthorized)
			return errors.Wrap(err, ErrExpiredToken)
		} else if errors.Is(err, ErrTokenTooEarly) {
			reqRes.SetStatusCode(httpstatus.Unauthorized)
			return errors.Wrap(err, ErrTooEarly)
		} else {
			reqRes.SetStatusCode(httpstatus.InternalServerError)
			return errors.Wrap(err, ErrUserUnauthorized)
		}
	}

	self.AuthorizedUser = *claims

	return nil
}
