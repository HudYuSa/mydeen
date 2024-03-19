package middlewares

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestTimeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), timeout)
		defer cancel()

		// Replace the current context with the one that has a timeout
		ctx.Request = ctx.Request.WithContext(requestTimeoutCtx)

		// Continue processing the request
		ctx.Next()
	}
}

func DBTimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbTimeoutCtx, cancel := context.WithTimeout(ctx.Request.Context(), timeout)
		defer cancel()

		// Replace the current context with the one that has a timeout
		ctx.Set("dbTimeoutContext", dbTimeoutCtx)

		// Continue processing the request
		ctx.Next()
	}
}
