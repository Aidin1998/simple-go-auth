terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
  required_version = ">= 1.0"
}

provider "aws" {
  region = var.aws_region
}

# Cognito User Pool
resource "aws_cognito_user_pool" "user_pool" {
  name = var.cognito_user_pool_name
}

# Cognito User Pool Client
resource "aws_cognito_user_pool_client" "user_pool_client" {
  name         = var.cognito_app_name
  user_pool_id = aws_cognito_user_pool.user_pool.id
}

# RDS Postgres Instance
resource "aws_db_instance" "postgres" {
  identifier         = "simple-go-auth-db"
  engine             = "postgres"
  instance_class     = var.db_instance_class
  allocated_storage  = var.db_allocated_storage
  name               = var.db_name
  username           = var.db_username
  password           = var.db_password
  skip_final_snapshot = true
}

# Secrets Manager for DB Password
resource "aws_secretsmanager_secret" "db_password_secret" {
  name = "db-password-secret"
}

resource "aws_secretsmanager_secret_version" "db_password_secret_version" {
  secret_id     = aws_secretsmanager_secret.db_password_secret.id
  secret_string = var.db_password
}

# Secrets Manager for Cognito App Secret
resource "aws_secretsmanager_secret" "cognito_app_secret" {
  name = "cognito-app-secret"
}

resource "aws_secretsmanager_secret_version" "cognito_app_secret_version" {
  secret_id     = aws_secretsmanager_secret.cognito_app_secret.id
  secret_string = aws_cognito_user_pool_client.user_pool_client.client_secret
}
