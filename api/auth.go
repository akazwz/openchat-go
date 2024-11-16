package api

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"golang.org/x/crypto/bcrypt"

	"openchat/dto"
	"openchat/lib"
	"openchat/model"
	"openchat/utils/rcontext"
	"openchat/utils/xhttp"
)

var Auth = &auth{}

type auth struct{}

func (auth) Signup(w http.ResponseWriter, r *http.Request) {
	var reqData dto.SignupReqData
	if err := xhttp.Bind(r, &reqData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var count int64
	if err := lib.DB.Model(model.User{}).Where("username = ?", reqData.Username).Limit(1).Count(&count).Error; err == nil && count == 1 {
		http.Error(w, "username already exists", http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqData.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user := model.User{
		Username:       reqData.Username,
		HashedPassword: string(hashedPassword),
	}
	if err := lib.DB.Create(&user).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	xhttp.RespJson(w, user, http.StatusCreated)
}

func (auth) Signin(w http.ResponseWriter, r *http.Request) {
	var reqData dto.SigninReqData
	if err := xhttp.Bind(r, &reqData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var user model.User
	if err := lib.DB.Model(model.User{}).Where("username = ?", reqData.Username).Limit(1).First(&user).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(reqData.Password)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	claims := &dto.MyClaims{
		UserId: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte("secret"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	randomId, err := gonanoid.New()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	refreshToken := model.RefreshToken{
		UserId:    user.ID,
		Token:     randomId,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
	}
	if err := lib.DB.Create(&refreshToken).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := dto.SigninRespData{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
	}
	xhttp.RespJson(w, data, http.StatusOK)
}

func (auth) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var reqData dto.RefreshTokenReqData
	if err := xhttp.Bind(r, &reqData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var refreshToken model.RefreshToken
	if err := lib.DB.Model(model.RefreshToken{}).Where("token = ?", reqData.RefreshToken).Limit(1).First(&refreshToken).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if refreshToken.ExpiresAt.Before(time.Now()) {
		http.Error(w, "refresh token expired", http.StatusBadRequest)
		return
	}
	claims := &dto.MyClaims{
		UserId: refreshToken.UserId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte("secret"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	xhttp.RespJson(w, dto.SigninRespData{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
	}, http.StatusOK)
}

func (auth) Account(w http.ResponseWriter, r *http.Request) {
	userId := rcontext.GetUserId(r.Context())
	var user model.User
	if err := lib.DB.Where("id = ?", userId).Limit(1).First(&user).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	xhttp.RespJson(w, user, http.StatusOK)
}
