package domain

import "github.com/golang-jwt/jwt"

var JwtKey = []byte("secret-key")

type JwtCustomClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}
