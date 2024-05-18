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
		// TODO: answer here
		sessionToken, err := ctx.Cookie("session_token")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		tokenClaims := &domain.JwtCustomClaims{}
		token, err := jwt.ParseWithClaims(sessionToken, tokenClaims, func(token *jwt.Token) (any, error) {
			return domain.JwtKey, nil
		})

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, web.NewBadRequestResponse("parse token failed"))
			return
		}

		if !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, web.NewUnauthorized("token is invalid"))
			return
		}

		ctx.Set("username", tokenClaims.Username)
		ctx.Next()
	})
}
