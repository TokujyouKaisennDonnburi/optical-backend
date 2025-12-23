variable "aws_region" {
  type    = string
  default = "ap-northeast-1"
}

variable "optical_domain" {
  type      = string
}

variable "jwt_secret" {
  type      = string
  sensitive = true
}

variable "redis_encryption_key" {
  type      = string
  sensitive = true
}

variable "agent_api_key" {
  type      = string
  sensitive = true
}

variable "github_client_id" {
  type = string
}

variable "github_client_secret" {
  type      = string
  sensitive = true
}

variable "github_oauth_redirect_uri" {
  type = string
}
