resource "aws_instance" "optical_backend_api" {
  ami           = "resolve:ssm:/aws/service/ami-amazon-linux-latest/al2023-ami-kernel-default-x86_64"
  instance_type = "t3.micro"
  tags = {
    Name = "optical_backend_api"
  }

  iam_instance_profile   = aws_iam_instance_profile.optical_api.name
  vpc_security_group_ids = [aws_security_group.optical_api_sg.id]
  # subnet_id              = aws_subnet.private_a.id
  subnet_id              = aws_subnet.public_sub_a.id

  user_data = <<EOF
#!/bin/bash
yum update -y
yum install -y ca-certificates
update-ca-trust force-enable
update-ca-trust extract
yum install -y docker jq
systemctl enable --now docker
usermod -a -G docker ec2-user
newgrp docker

docker pull ghcr.io/tokujyoukaisenndonnburi/optical-backend:latest

cat <<'SCRIPT67' > /home/ec2-user/entrypoint.sh
#!/bin/bash
PSQL_SECRET_ARN="${aws_db_instance.optical_psql.master_user_secret[0].secret_arn}"
PSQL_SECRET_JSON=$(aws secretsmanager get-secret-value \
  --secret-id "$PSQL_SECRET_ARN" \
  --query SecretString \
  --output text)
APP_SECRET_JSON=$(aws secretsmanager get-secret-value \
  --secret-id ${aws_secretsmanager_secret.optical_app_secrets.name} \
  --query SecretString \
  --output text)

export POSTGRES_HOST=${aws_db_instance.optical_psql.address}
export POSTGRES_PORT=${aws_db_instance.optical_psql.port}
export POSTGRES_DB=${aws_db_instance.optical_psql.db_name}
export POSTGRES_USER=$(echo "$PSQL_SECRET_JSON" | jq -r .username)
export POSTGRES_PASSWORD=$(echo "$PSQL_SECRET_JSON" | jq -r .password)
export POSTGRES_SSLMODE=require

export REDIS_TLS=1
export REDIS_ADDRESS="${aws_elasticache_replication_group.optical_redis.primary_endpoint_address}:${aws_elasticache_replication_group.optical_redis.port}"
export REDIS_PASSWORD=$(aws ssm get-parameter --name "${aws_ssm_parameter.redis_password.name}" --with-decryption --query Parameter.Value --output text)

export MINIO_FROM_IAM=1
export MINIO_REGION="${var.aws_region}"
export MINIO_ENDPOINT="s3.${var.aws_region}.amazonaws.com"
export MINIO_IMAGE_BUCKET_NAME="${aws_s3_bucket.optical_images_bucket.bucket}"
export IMAGE_STORAGE_BASE_URL="https://${aws_s3_bucket.optical_images_bucket.bucket}.s3.${var.aws_region}.amazonaws.com"

export JWT_SECRET_KEY=$(echo "$APP_SECRET_JSON" | jq -r .JWT_SECRET)
export REDIS_ENCRYPTION_KEY=$(echo "$APP_SECRET_JSON" | jq -r .REDIS_ENCRYPTION_KEY)

export AGENT_MODEL="openai/gpt-oss-120b"
export AGENT_API_KEY=$(echo "$APP_SECRET_JSON" | jq -r .AGENT_API_KEY)

export GITHUB_APP_ID="${var.github_client_id}"
export GITHUB_CLIENT_ID="${var.github_client_id}"
export GITHUB_CLIENT_SECRET=$(echo "$APP_SECRET_JSON" | jq -r .GITHUB_CLIENT_SECRET)
export GITHUB_OAUTH_REDIRECT_URI="${var.github_oauth_redirect_uri}"

export RUNTIME_MIGRATION=1
export LOGGER_JSON_FORMAT=1
export FRONTEND_BASE_URL="http://localhost:3000"

export EMAIL_HOST="smtp.gmail.com"
export EMAIL_PORT=587
export EMAIL_USER=""
export EMAIL_PASSWORD=""

docker run --rm \
  -e POSTGRES_HOST \
  -e POSTGRES_PORT \
  -e POSTGRES_DB \
  -e POSTGRES_USER \
  -e POSTGRES_PASSWORD \
  -e POSTGRES_SSLMODE \
  -e REDIS_TLS \
  -e REDIS_ADDRESS \
  -e REDIS_PASSWORD \
  -e MINIO_FROM_IAM \
  -e MINIO_REGION \
  -e MINIO_ENDPOINT \
  -e MINIO_IMAGE_BUCKET_NAME \
  -e IMAGE_STORAGE_BASE_URL \
  -e JWT_SECRET_KEY \
  -e REDIS_ENCRYPTION_KEY \
  -e AGENT_MODEL \
  -e AGENT_API_KEY \
  -e GITHUB_CLIENT_ID \
  -e GITHUB_CLIENT_SECRET \
  -e GITHUB_OAUTH_REDIRECT_URI \
  -e RUNTIME_MIGRATION \
  -e LOGGER_JSON_FORMAT \
  -e FRONTEND_BASE_URL \
  -e EMAIL_HOST \
  -e EMAIL_PORT \
  -e EMAIL_USER \
  -e EMAIL_PASSWORD \
  --name optical_backend \
  -p 8000:8000 \
  ghcr.io/tokujyoukaisenndonnburi/optical-backend:latest
SCRIPT67

chmod +x /home/ec2-user/entrypoint.sh
chown ec2-user:ec2-user /home/ec2-user/entrypoint.sh

sudo tee /etc/systemd/system/optical-backend.service > /dev/null <<'UNIT67'
[Unit]
Description=Optical Backend API Server
After=network.target

[Service]
User=ec2-user
Group=ec2-user
WorkingDirectory=/home/ec2-user

ExecStart=/home/ec2-user/entrypoint.sh
Restart=always

[Install]
WantedBy=multi-user.target
UNIT67

sudo systemctl daemon-reload
sudo systemctl enable optical-backend.service
sudo systemctl start optical-backend.service

EOF
}
