// auth/auth_handler.go
package auth

import (
	"regexp"

	"my-go-project/config"

	"github.com/labstack/echo/v4"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	Service         *AuthServiceImpl
	RecaptchaSecret string
}

// NewHandler registers all auth routes on the given Echo and returns the handler.
// You must pass in your app-wide config so we can read RecaptchaSecretKey.
func NewHandler(e *echo.Echo, svc *AuthServiceImpl, cfg *config.Config) *AuthHandler {
	h := &AuthHandler{
		Service:         svc,
		RecaptchaSecret: cfg.RecaptchaSecretKey,
	}

	// Public routes
	e.GET("/ping", h.Ping)
	e.POST("/signup", h.SignUp)
	e.POST("/signin", h.SignIn)
	e.POST("/confirm", h.ConfirmSignUp)
	e.POST("/refresh", h.Refresh)

	// Protected
	e.POST("/logout", h.SignOut, NewMiddleware(svc))

	return h
}

// Ping returns pong.
func (h *AuthHandler) Ping(c echo.Context) error {
	return c.String(200, "pong")
}

// SignUp handles user registration + captcha + password policy.
func (h *AuthHandler) SignUp(c echo.Context) error {
	var req struct {
		Username     string `json:"username"`
		Email        string `json:"email"`
		Password     string `json:"password"`
		CaptchaToken string `json:"captcha_token"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "invalid request body"})
	}

	// 1) Verify reCAPTCHA
	if !verifyCaptcha(req.CaptchaToken, h.RecaptchaSecret) {
		return c.JSON(400, map[string]string{"error": "captcha failed"})
	}

	// 2) Enforce password policy
	pwdPolicy := regexp.MustCompile(`^(?=.*[0-9])(?=.*[A-Z]).{8,}$`)
	if !pwdPolicy.MatchString(req.Password) {
		return c.JSON(400, map[string]string{"error": "password must be ≥8 chars, include uppercase & digit"})
	}

	// 3) Call service.SignUp
	if err := h.Service.SignUp(c.Request().Context(), req.Username, req.Password, req.Email); err != nil {
		return c.JSON(500, map[string]string{"error": "failed to create user"})
	}

	// 4) Return 202—Cognito has sent the confirmation code
	return c.JSON(202, map[string]string{"message": "verification code sent"})
}

// verifyCaptcha validates the reCAPTCHA token using the secret key.
// Currently a stub; replace with real HTTP call in prod.
func verifyCaptcha(token, secret string) bool {
	return true
}

// ConfirmSignUp confirms a user's sign-up using a verification code.
func (h *AuthHandler) ConfirmSignUp(c echo.Context) error {
	var req struct {
		Username string `json:"username"`
		Code     string `json:"code"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "invalid request body"})
	}
	if err := h.Service.ConfirmSignUp(c.Request().Context(), req.Username, req.Code); err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}
	return c.JSON(200, map[string]string{"message": "user confirmed"})
}

// SignIn authenticates and returns JWT + refresh token.
func (h *AuthHandler) SignIn(c echo.Context) error {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "invalid request body"})
	}
	tokens, err := h.Service.SignIn(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		return c.JSON(401, map[string]string{"error": "invalid credentials"})
	}
	return c.JSON(200, tokens)
}

// Refresh rotates tokens.
func (h *AuthHandler) Refresh(c echo.Context) error {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "invalid request body"})
	}
	tokens, err := h.Service.RefreshTokens(c.Request().Context(), req.RefreshToken)
	if err != nil {
		return c.JSON(401, map[string]string{"error": err.Error()})
	}
	return c.JSON(200, tokens)
}

// SignOut revokes the user’s session.
func (h *AuthHandler) SignOut(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.JSON(400, map[string]string{"error": "authorization token required"})
	}
	if err := h.Service.SignOut(c.Request().Context(), token); err != nil {
		return c.JSON(500, map[string]string{"error": "failed to sign out"})
	}
	return c.NoContent(204)
}
