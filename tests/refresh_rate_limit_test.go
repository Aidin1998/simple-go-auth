// tests/refresh_rate_limit_test.go
package tests

import (
	"net/http" // ← standard http for StatusTooManyRequests, StatusOK, etc.
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	// Removed unused import for "my-go-project/http"
)

func TestRateLimitOnSignin(t *testing.T) {
	// 1) Spin up a fresh Echo instance
	e := echo.New()

	// 2) Attach a dummy /signin endpoint with a 10-req/min rate limiter
	rlStore := middleware.NewRateLimiterMemoryStore(10)
	e.POST("/signin",
		// dummy handler always returns 200
		func(c echo.Context) error {
			return c.NoContent(http.StatusOK)
		},
		// rate-limit middleware
		middleware.RateLimiter(rlStore),
	)

	// 3) Send 11 requests; first 10 should NOT be 429, the 11th MUST be 429
	for i := 0; i < 11; i++ {
		req := httptest.NewRequest("POST", "/signin", strings.NewReader(`{"username":"x","password":"y"}`))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		if i < 10 {
			assert.NotEqual(t, http.StatusTooManyRequests, rec.Code,
				"request #%d should NOT be rate-limited", i+1)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, rec.Code,
				"request #%d should be rate-limited", i+1)
		}
	}
}
