package pkg

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_SignupSigninRefreshProtected(t *testing.T) {
	t := assert.New(t)

	// Mock server setup
	handler := setupRouter() // Assuming setupRouter initializes routes
	srv := httptest.NewServer(handler)
	defer srv.Close()

	// Test Signup
	signupPayload := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	resp, err := http.Post(srv.URL+"/signup", "application/json", bytes.NewBuffer(json.Marshal(signupPayload)))
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	// Test Signin
	signinPayload := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	resp, err := http.Post(srv.URL+"/signin", "application/json", bytes.NewBuffer(json.Marshal(signupPayload)))
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	// Test Refresh Token
	refreshPayload := map[string]string{
		"refresh_token": "mock_refresh_token", // Replace with actual token from signin response
	}
	refreshReqBody, _ := json.Marshal(refreshPayload)
	refreshReq, _ := http.NewRequest("POST", srv.URL+"/refresh", bytes.NewBuffer(refreshReqBody))
	refreshReq.Header.Set("Content-Type", "application/json")
	refreshResp, err := http.DefaultClient.Do(refreshReq)
	t.NoError(err)
	t.Equal(http.StatusOK, refreshResp.StatusCode)

	// Test Protected Endpoint
	protectedReq, _ := http.NewRequest("GET", srv.URL+"/protected", nil)
	protectedReq.Header.Set("Authorization", "Bearer mock_access_token") // Replace with actual token
	protectedResp, err := http.DefaultClient.Do(protectedReq)
	t.NoError(err)
	t.Equal(http.StatusOK, protectedResp.StatusCode)
}