package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *AuthHandler) InitiateMFASetup(c echo.Context) error {
	// Call CognitoClient.AssociateSoftwareToken
	return c.JSON(http.StatusOK, map[string]string{"message": "MFA setup initiated"})
}

func (h *AuthHandler) VerifyMFAToken(c echo.Context) error {
	// Call CognitoClient.VerifySoftwareToken
	return c.JSON(http.StatusOK, map[string]string{"message": "MFA verified"})
}
