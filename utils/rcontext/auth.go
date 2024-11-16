package rcontext

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"openchat/dto"
)

func WithUserId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" || len(authorization) <= len("Bearer ") {
			http.Error(w, "authorization header is required", http.StatusUnauthorized)
			return
		}
		authorization = authorization[len("Bearer "):]
		token, err := jwt.ParseWithClaims(authorization, &dto.MyClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(*dto.MyClaims)
		if !ok {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user_id", claims.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserId(ctx context.Context) string {
	userId, ok := ctx.Value("user_id").(string)
	if !ok {
		return ""
	}
	return userId
}
