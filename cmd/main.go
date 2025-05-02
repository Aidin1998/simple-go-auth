// cmd/main.go
package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.uber.org/zap"

	"my-go-project/auth"
	"my-go-project/aws"
	"my-go-project/config"
	"my-go-project/db"
	"my-go-project/http"
	"my-go-project/otel"
)

func main() {
	// 0. Initialize OpenTelemetry
	shutdown := otel.InitTracer()
	defer shutdown()

	// 1) Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2) Init DB
	dbInstance, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 3) Init Cognito client
	cognitoClient, err := aws.NewCognitoClient(
		cfg.AWSRegion,
		cfg.CognitoUserPoolID,
		cfg.CognitoAppClientID,
	)
	if err != nil {
		log.Fatalf("Failed to initialize Cognito client: %v", err)
	}

	// 4) Build AuthService
	authService := auth.NewAuthServiceImpl(cognitoClient, dbInstance)

	// 5) Build your handler (this also registers its own routes on a new echo.Group internally)
	//    Note: it DOES NOT create the base echo - just records handler methods.
	authHandler := auth.NewHandler(nil, authService, cfg) // we’ll pass 'nil' because SetupRouter will mount routes directly

	// Create shared Zap logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize Zap logger: %v", err)
	}
	defer logger.Sync()

	// 3bis. Wrap Echo with OTel middleware
	e := echo.New()
	e.Use(otelecho.Middleware("my-go-auth-service"))

	// 6) Create the Echo router with global middleware + config/ping + auth routes
	router := http.SetupRouter(authHandler, auth.NewMiddleware(authService), logger)

	// Set Echo server read and write timeouts
	router.Server.ReadTimeout = cfg.EchoReadTimeout
	router.Server.WriteTimeout = cfg.EchoWriteTimeout

	// 7) Root welcome (optional – you can also add this in SetupRouter)
	router.GET("/", func(c echo.Context) error {
		return c.String(200, "🚀 Welcome to BitPolaris! Please POST to /signin or /signup")
	})

	// 8) Start
	log.Printf("Server running on port %s...\n", cfg.Port)

	// Listen with TLS (HTTP/2)
	log.Fatal(router.StartTLS(":"+cfg.Port, "certs/server.crt", "certs/server.key"))
}
