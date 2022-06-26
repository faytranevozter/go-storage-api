package domain

import (
	"net/http"
	"net/url"

	"github.com/golang-jwt/jwt/v4"
)

// default param
type DefaultPayload struct {
	Query    url.Values
	Request  *http.Request
	ID       interface{}
	Payload  interface{}
	AuthData *JWTPayload
}

type JWTPayload struct {
	ID int64 `json:"id"`
	jwt.StandardClaims
}
