package helpers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		logrus.WithFields(logrus.Fields{
			"ip":      c.ClientIP(),
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
			"proto":   c.Request.Proto,
			"status":  c.Writer.Status(),
			"latency": time.Since(startTime),
			"ua":      c.Request.UserAgent(),
		}).Info()
	}
}

func CustomRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logrus.Error("Panic Recover : ", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, ErrResp(http.StatusInternalServerError, "Something went wrong"))
			}
		}()
		c.Next()
	}
}
