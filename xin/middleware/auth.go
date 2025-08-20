package middleware

import (
	"github.com/askasoft/pango/xin"
)

const (
	AuthUserKey = "X_AUTH_USER" // Key for authenticated user object saved in context
)

// AuthUser a authenticated user interface
type AuthUser interface {
	GetUsername() string
	GetPassword() string
}

// FindUserFunc user lookup function
type FindUserFunc func(c *xin.Context, username, password string) (AuthUser, error)
