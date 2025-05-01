package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPingEndpoint(t *testing.T) {
	// Initialize Echo instance
	e := echo.New()

	// Define the /ping route
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	// Create a fake GET /ping request
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Serve the request
	e.Router().Find(http.MethodGet, "/ping", c)
	e.ServeHTTP(rec, req)

	// Assert we got a 200 and body "pong"
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "pong", rec.Body.String())
}
