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

		if err != nil || sessionToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, web.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Status:  "status unauthorized",
				Message: "session token is empty",
			})
			return
		}

		tokenClaims := &domain.JwtCustomClaims{}
		token, err := jwt.ParseWithClaims(sessionToken, tokenClaims, func(token *jwt.Token) (any, error) {
			return domain.JwtKey, err
		})

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, web.ErrorResponse{
				Code:    http.StatusBadRequest,
				Status:  "status bad request",
				Message: "parse token failed",
			})
			return
		}

		if !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, web.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Status:  "status unauthorized",
				Message: "token is invalid",
			})
			return
		}

		ctx.Set("username", tokenClaims.Username)
		ctx.Next()
	})
}
