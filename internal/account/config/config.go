package config

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/spf13/viper"
)

type Config struct {
	Port      string
	DBDSN     string
	JWTSecret string
	RedisURL  string
}

var Cfg Config

func Init() {
	viper.SetEnvPrefix("ACCOUNT")
	viper.AutomaticEnv()
	Cfg.Port = viper.GetString("PORT")
	if Cfg.Port == "" {
		Cfg.Port = "8080"
	}
	Cfg.RedisURL = viper.GetString("REDIS_URL")
	if Cfg.RedisURL == "" {
		log.Fatal("REDIS_URL must be set")
	}

	// Load AWS Secrets Manager
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("AWS config load error: %v", err)
	}
	sm := secretsmanager.NewFromConfig(awsCfg)
	secOut, err := sm.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String("account-service-secrets"),
	})
	if err != nil {
		log.Fatalf("unable to fetch secrets: %v", err)
	}
	var s struct {
		DBDSN     string `json:"db_dsn"`
		JWTSecret string `json:"jwt_secret"`
	}
	if err := json.Unmarshal([]byte(*secOut.SecretString), &s); err != nil {
		log.Fatalf("invalid secret JSON: %v", err)
	}
	Cfg.DBDSN = s.DBDSN
	Cfg.JWTSecret = s.JWTSecret
}
