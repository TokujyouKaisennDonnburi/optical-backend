resource "aws_secretsmanager_secret" "optical_app_secrets" {
  name = "optical/app/secrets"
  description = "enviroment variables for optical_backend"
  recovery_window_in_days = 0
}

resource "aws_secretsmanager_secret_version" "optical_app_secrets" {
  secret_id = aws_secretsmanager_secret.optical_app_secrets.id

  secret_string = jsonencode({
    JWT_SECRET           = var.jwt_secret
    REDIS_ENCRYPTION_KEY = var.redis_encryption_key
    AGENT_API_KEY        = var.agent_api_key
    GITHUB_CLIENT_ID     = var.github_client_id
    GITHUB_CLIENT_SECRET = var.github_client_secret
  })
}
