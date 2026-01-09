resource "random_password" "optical_redis" {
  length           = 64
  special          = false
}

resource "aws_elasticache_replication_group" "optical_redis" {
  replication_group_id  = "optical-redis"
  description = "Redis"
  engine         = "redis"
  engine_version = "7.1"
  node_type      = "cache.t4g.small"
  port           = 6379
  num_cache_clusters = 1

  subnet_group_name  = aws_elasticache_subnet_group.optical_redis.name
  security_group_ids = [aws_security_group.optical_redis_sg.id]

  transit_encryption_enabled = true
  auth_token                 = random_password.optical_redis.result
}

resource "aws_elasticache_subnet_group" "optical_redis" {
  name = "optical-redis-subnet"
  subnet_ids = [
    aws_subnet.private_a.id,
    aws_subnet.private_c.id,
  ]
}

resource "aws_ssm_parameter" "redis_password" {
  name  = "/optical_api/optical_redis/password"
  type  = "SecureString"
  value = random_password.optical_redis.result
}
