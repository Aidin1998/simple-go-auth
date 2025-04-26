# Simple Go Auth

## Overview
This project provides a simple authentication service using Go, JWT, and AWS Secrets Manager. It is designed to run on AWS ECS behind an Application Load Balancer and scale to handle up to 100,000 concurrent users.

## Local Running Instructions
1. Clone the repository:
   ```bash
   git clone https://github.com/your-repo/simple-go-auth.git
   cd simple-go-auth/my-go-project
   ```

2. Build and run the application locally:
   ```bash
   go run ./cmd/main.go
   ```

3. Access the application at `http://localhost:80`.

## Local Testing Instructions
1. Run unit tests:
   ```bash
   go test ./tests -cover
   ```

2. Run integration tests:
   ```bash
   go test ./tests/auth_integration_test.go
   ```

## Deployment Instructions
1. Build the Docker image:
   ```bash
   docker build -t simple-go-auth .
   ```

2. Push the Docker image to AWS ECR:
   ```bash
   aws ecr get-login-password --region <region> | docker login --username AWS --password-stdin <account_id>.dkr.ecr.<region>.amazonaws.com
   docker tag simple-go-auth:latest <account_id>.dkr.ecr.<region>.amazonaws.com/simple-go-auth:latest
   docker push <account_id>.dkr.ecr.<region>.amazonaws.com/simple-go-auth:latest
   ```

3. Update the ECS service:
   ```bash
   aws ecs update-service --cluster <cluster_name> --service <service_name> --force-new-deployment
   ```

## AWS Secrets Setup Instructions
1. Create a secret in AWS Secrets Manager with the key `jwtSecretKey` and your desired secret value.
2. Update the ECS task definition to include the secret as an environment variable.

## Notes
- Ensure that your AWS IAM roles have the necessary permissions for ECS, ECR, and Secrets Manager.
- The application is optimized for production with a multi-stage Docker build and CI/CD pipeline.