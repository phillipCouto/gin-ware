package timeout

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// Timeout replaces the context in the request with a new one that
// has a timeout that matches the provided duration. Internally the
// middleware is relying on the context to properly cancel any in progress
// work so all code further down the pipeline should be context aware.
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.Clone(ctx)
		c.Next()
	}
}
