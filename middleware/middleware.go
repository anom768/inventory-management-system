package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"inventory-management-system/model/domain"
	"inventory-management-system/model/web"
	"net/http"
)

func Auth() gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		sessionToken, err := ctx.Cookie("session_token")
		if err != nil || sessionToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, web.NewUnauthorizedError("session token is empty"))
			return
		}

		tokenClaims := &domain.JwtCustomClaims{}
		token, err := jwt.ParseWithClaims(sessionToken, tokenClaims, func(token *jwt.Token) (any, error) {
			return domain.JwtKey, err
		})

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestError("parse token failed"))
			return
		}

		if !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, web.NewUnauthorizedError("token is invalid"))
			return
		}
		
		ctx.Set("username", tokenClaims.Username)
		ctx.Set("role", tokenClaims.Role)
		ctx.Next()
	})
}

func AdminOnly() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("role")
		if !exists || role != "admin" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, web.NewUnauthorizedError("user is not admin"))
			return
		}

		ctx.Next()
	}
}
