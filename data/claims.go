package data

import jwt "github.com/dgrijalva/jwt-go"

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}
