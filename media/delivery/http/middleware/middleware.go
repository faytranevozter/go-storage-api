package middleware

import (
	"net/http"
	"os"
	"strings"

	"storage-api/app/helpers"
	"storage-api/domain"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// GoMiddleware represent the data-struct for middleware
type GoMiddleware struct {
	secretKey string
}

// Auth for jwt authentication
func (m *GoMiddleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get header
		hAuth := c.GetHeader("Authorization")

		splitToken := strings.Split(hAuth, "Bearer ")
		if len(splitToken) != 2 {
			c.JSON(http.StatusUnauthorized, helpers.ErrResp(http.StatusUnauthorized, "Unauthorized"))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenString := splitToken[1]

		token, err := jwt.ParseWithClaims(tokenString, &domain.JWTPayload{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.secretKey), nil
		})

		if token == nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, helpers.ErrResp(http.StatusUnauthorized, err.Error()))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*domain.JWTPayload)
		if !ok {
			c.JSON(http.StatusUnauthorized, helpers.ErrResp(http.StatusUnauthorized, "invalid token data"))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("JWTDATA", claims)

		c.Next()
	}
}

// InitMiddleware initialize the middleware
func InitMiddleware() *GoMiddleware {
	return &GoMiddleware{
		secretKey: os.Getenv("JWT_SECRET_KEY"),
	}
}
