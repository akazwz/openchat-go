package dto

import "github.com/golang-jwt/jwt/v5"

type SignupReqData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SigninReqData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshTokenReqData struct {
	RefreshToken string `json:"refresh_token"`
}

type SigninRespData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type MyClaims struct {
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}
