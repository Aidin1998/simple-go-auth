output "cognito_user_pool_id" {
  description = "ID of the Cognito User Pool"
  value       = aws_cognito_user_pool.user_pool.id
}

output "cognito_app_client_id" {
  description = "ID of the Cognito App Client"
  value       = aws_cognito_user_pool_client.user_pool_client.id
}

output "rds_endpoint" {
  description = "Endpoint of the RDS Postgres instance"
  value       = aws_db_instance.postgres.endpoint
}

output "db_password_secret_arn" {
  description = "ARN of the Secrets Manager secret for the DB password"
  value       = aws_secretsmanager_secret.db_password_secret.arn
}
