variable "aws_region" {
  description = "AWS region to deploy resources"
  type        = string
}

variable "cognito_user_pool_name" {
  description = "Name of the Cognito User Pool"
  type        = string
}

variable "cognito_app_name" {
  description = "Name of the Cognito App Client"
  type        = string
}

variable "db_username" {
  description = "Database username"
  type        = string
}

variable "db_password" {
  description = "Database password"
  type        = string
  sensitive   = true
}

variable "db_name" {
  description = "Database name"
  type        = string
}

variable "db_instance_class" {
  description = "Database instance class"
  type        = string
}

variable "db_allocated_storage" {
  description = "Allocated storage for the database (in GB)"
  type        = number
}
