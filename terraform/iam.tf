resource "aws_iam_instance_profile" "optical_api" {
  role = aws_iam_role.optical_api.name
}

resource "aws_iam_role" "optical_api" {
  name = "optical_api_secrets_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = { Service = "ec2.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_role_policy_attachment" "secrets_attach" {
  role       = aws_iam_role.optical_api.name
  policy_arn = aws_iam_policy.secrets_read.arn
}

resource "aws_iam_role_policy_attachment" "s3_attach" {
  role       = aws_iam_role.optical_api.name
  policy_arn = aws_iam_policy.s3_access.arn
}

resource "aws_iam_role_policy_attachment" "redis_attach" {
  role       = aws_iam_role.optical_api.name
  policy_arn = aws_iam_policy.read_redis_pass.arn
}

resource "aws_iam_role_policy_attachment" "app_secrests_attach" {
  role       = aws_iam_role.optical_api.name
  policy_arn = aws_iam_policy.app_secrets.arn
}

resource "aws_iam_policy" "secrets_read" {
  name = "secretsmanager-read"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "secretsmanager:GetSecretValue"
      ]
      Resource = aws_db_instance.optical_psql.master_user_secret[0].secret_arn
    }]
  })
}

resource "aws_iam_policy" "s3_access" {
  name = "s3-access"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "s3:GetObject",
        "s3:PutObject",
        "s3:ListBucket"
      ]
      Resource = [
        aws_s3_bucket.optical_images_bucket.arn,
        "${aws_s3_bucket.optical_images_bucket.arn}/*"
      ]
    }]
  })
}

resource "aws_iam_policy" "read_redis_pass" {
  name = "read_redis_password"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "ssm:GetParameter",
          "ssm:GetParameters"
        ]
        Resource = aws_ssm_parameter.redis_password.arn
      }
    ]
  })
}

resource "aws_iam_policy" "app_secrets" {
  name = "app-secrets-read"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "secretsmanager:GetSecretValue"
      ]
      Resource = aws_secretsmanager_secret.optical_app_secrets.arn
    }]
  })
}

resource "aws_iam_role_policy_attachment" "ec2_ssm" {
  role      = aws_iam_role.optical_api.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}
